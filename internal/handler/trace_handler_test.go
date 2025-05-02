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
	"github.com/zkropotkine/agent-trace/internal/model"
	"github.com/zkropotkine/agent-trace/internal/repository"
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

func TestPostTrace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		input          interface{}
		repoReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid trace",
			input:          model.Trace{TraceID: "123", AgentName: "AgentX"},
			repoReturn:     nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"trace saved"}`,
		},
		{
			name:           "Invalid JSON",
			input:          "{invalid_json}",
			repoReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid trace payload"}`,
		},
		{
			name:           "Repo error",
			input:          model.Trace{TraceID: "fail", AgentName: "AgentX"},
			repoReturn:     errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to save trace"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTraceRepo)
			h := NewTraceHandler(repo)

			r := gin.Default()
			r.POST("/trace", h.PostTrace)

			var req *http.Request

			switch input := tt.input.(type) {
			case string: // raw invalid JSON
				req = httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBufferString(input))
			default:
				body, _ := json.Marshal(input)
				req = httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBuffer(body))
				if trace, ok := input.(model.Trace); ok && tt.repoReturn != nil {
					repo.On("InsertTrace", mock.Anything, mock.MatchedBy(func(t model.Trace) bool {
						return t.TraceID == trace.TraceID && t.AgentName == trace.AgentName
					})).Return(tt.repoReturn).Once()
				} else if trace, ok := input.(model.Trace); ok {
					repo.On("InsertTrace", mock.Anything, mock.MatchedBy(func(t model.Trace) bool {
						return t.TraceID == trace.TraceID && t.AgentName == trace.AgentName
					})).Return(nil).Once()
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			repo.AssertExpectations(t)
		})
	}
}

func TestGetTracesHandler(t *testing.T) {
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
			name: "returns filtered traces",
			path: "/api/traces?agent=test-agent",
			setupMock: func(repo *mockTraceRepo) {
				repo.On("GetTraces", mock.Anything, mock.MatchedBy(func(f repository.TraceFilter) bool {
					return f.AgentName == "test-agent"
				})).Return([]model.Trace{{TraceID: "1", AgentName: "test-agent", Timestamp: now}}, nil)
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, body []byte) {
				var result map[string][]model.Trace
				err := json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Len(t, result["traces"], 1)
				assert.Equal(t, "test-agent", result["traces"][0].AgentName)
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
			r.GET("/api/traces", h.GetTraces)

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
