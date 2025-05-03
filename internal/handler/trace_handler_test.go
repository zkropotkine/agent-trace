package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zkropotkine/agent-trace/config"
	"github.com/zkropotkine/agent-trace/internal/model"
	"github.com/zkropotkine/agent-trace/internal/repository"
	"github.com/zkropotkine/agent-trace/pkg/logger"
)

type mockTraceRepo struct {
	mock.Mock
}

func (m *mockTraceRepo) InsertTrace(ctx context.Context, trace model.Trace) error {
	args := m.Called(ctx, trace)
	return args.Error(0)
}

func (m *mockTraceRepo) GetTraces(ctx context.Context, filter repository.TraceFilter) ([]model.Trace, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]model.Trace), args.Error(1)
}

func (m *mockTraceRepo) GetByID(ctx context.Context, id string) (*model.Trace, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Trace), args.Error(1)
}

func TestPostTraceRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		input          model.Trace
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid trace",
			input: model.Trace{
				TraceID:      "xyz789",
				AgentName:    "TestAgent",
				InputPrompt:  "Hello",
				OutputPrompt: "Hi there",
				Model:        "gpt-4-turbo",
			},
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"trace saved"}`,
		},
		{
			name: "repo error",
			input: model.Trace{
				TraceID:      "fail",
				AgentName:    "TestAgent",
				InputPrompt:  "Hello",
				OutputPrompt: "Hi there",
				Model:        "gpt-4-turbo",
			},
			mockReturn:     errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to save trace"}`,
		},
		{
			name: "repo error",
			input: model.Trace{
				TraceID:      "fail",
				AgentName:    "TestAgent",
				InputPrompt:  "Hello",
				OutputPrompt: "Hi there",
				Model:        "gpt-4-turbo-nonexistent",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to analyze tokens"}`,
			mockReturn:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTraceRepo)

			// TODO improve this code
			if !(tt.mockReturn == nil && tt.expectedStatus == http.StatusInternalServerError) {
				repo.On("InsertTrace", mock.Anything, mock.MatchedBy(func(tr model.Trace) bool {
					return tr.TraceID == tt.input.TraceID && tr.AgentName == tt.input.AgentName
				})).Return(tt.mockReturn)
			}

			h := NewTraceHandler(repo)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				ctx := logger.WithLogger(c.Request.Context(), logger.NewLogRusLogger(config.Log{}))
				c.Request = c.Request.WithContext(ctx)
				c.Next()
			})
			r.POST("/trace", h.PostTrace)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			repo.AssertExpectations(t)
		})
	}

	t.Run("invalid json", func(t *testing.T) {
		repo := new(mockTraceRepo)
		h := NewTraceHandler(repo)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			ctx := logger.WithLogger(c.Request.Context(), logger.NewLogRusLogger(config.Log{}))
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		})
		r.POST("/trace", h.PostTrace)

		req := httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBufferString("{invalid_json}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error":"invalid trace payload"}`, rec.Body.String())
		repo.AssertExpectations(t)
	})
}

func TestGetTracesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("query params from/to/limit/offset covered", func(t *testing.T) {
		repo := new(mockTraceRepo)
		h := NewTraceHandler(repo)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			ctx := logger.WithLogger(c.Request.Context(), logger.NewLogRusLogger(config.Log{}))
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		})
		r.GET("/traces", h.GetTraces)

		now := time.Now().UTC().Format(time.RFC3339)
		repo.On("GetTraces", mock.Anything, mock.MatchedBy(func(f repository.TraceFilter) bool {
			return f.AgentName == "agent1" && f.From != nil && f.To != nil && f.Limit == 10 && f.Offset == 5
		})).Return([]model.Trace{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/traces?agent=agent1&from="+now+"&to="+now+"&limit=10&offset=5", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"traces":[]}`, rec.Body.String())
		repo.AssertExpectations(t)
	})

	t.Run("repo error handling", func(t *testing.T) {
		repo := new(mockTraceRepo)
		h := NewTraceHandler(repo)
		r := gin.New()
		r.Use(func(c *gin.Context) {
			ctx := logger.WithLogger(c.Request.Context(), logger.NewLogRusLogger(config.Log{}))
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		})
		r.GET("/traces", h.GetTraces)

		repo.On("GetTraces", mock.Anything, mock.Anything).Return([]model.Trace(nil), errors.New("db fail"))
		req := httptest.NewRequest(http.MethodGet, "/traces", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"failed to fetch traces"}`, rec.Body.String())
	})

	t.Run("default limit and offset", func(t *testing.T) {
		repo := new(mockTraceRepo)
		h := NewTraceHandler(repo)
		r := gin.New()
		r.Use(func(c *gin.Context) {
			ctx := logger.WithLogger(c.Request.Context(), logger.NewLogRusLogger(config.Log{}))
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		})
		r.GET("/traces", h.GetTraces)

		repo.On("GetTraces", mock.Anything, mock.MatchedBy(func(f repository.TraceFilter) bool {
			return f.Limit == 50 && f.Offset == 0
		})).Return([]model.Trace{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/traces", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"traces":[]}`, rec.Body.String())
		repo.AssertExpectations(t)
	})
}

func TestGetTraceByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	tests := []struct {
		name           string
		path           string
		setupMock      func(repo *mockTraceRepo)
		expectedStatus int
		assertBody     func(t *testing.T, body []byte)
	}{
		{
			name: "returns trace by id",
			path: "/api/traces/abc123",
			setupMock: func(repo *mockTraceRepo) {
				repo.On("GetByID", mock.Anything, "abc123").Return(&model.Trace{TraceID: "abc123", AgentName: "test-agent", Timestamp: now}, nil)
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, body []byte) {
				var trace model.Trace
				err := json.Unmarshal(body, &trace)
				assert.NoError(t, err)
				assert.Equal(t, "abc123", trace.TraceID)
			},
		},
		{
			name: "returns 404 if trace not found",
			path: "/api/traces/missing",
			setupMock: func(repo *mockTraceRepo) {
				repo.On("GetByID", mock.Anything, "missing").Return((*model.Trace)(nil), errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			assertBody: func(t *testing.T, body []byte) {
				var res map[string]string
				err := json.Unmarshal(body, &res)
				assert.NoError(t, err)
				assert.Equal(t, "trace not found", res["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTraceRepo)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			h := NewTraceHandler(repo)
			r := gin.New()
			r.GET("/api/traces/:id", h.GetTraceByID)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.assertBody != nil {
				tt.assertBody(t, resp.Body.Bytes())
			}
			repo.AssertExpectations(t)
		})
	}
}
