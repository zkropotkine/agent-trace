package router

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
	"github.com/zkropotkine/agent-trace/internal/handler"
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
			h := handler.NewTraceHandler(repo)

			r := gin.Default()
			r.POST("/trace", h.PostTrace)

			var req *http.Request

			switch input := tt.input.(type) {
			case string:
				req = httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBufferString(input))
			default:
				body, _ := json.Marshal(input)
				req = httptest.NewRequest(http.MethodPost, "/trace", bytes.NewBuffer(body))
				if trace, ok := input.(model.Trace); ok {
					repo.On("InsertTrace", mock.Anything, mock.MatchedBy(func(t model.Trace) bool {
						return t.TraceID == trace.TraceID && t.AgentName == trace.AgentName
					})).Return(tt.repoReturn).Once()
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

func TestGetTracesRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()

	repo := new(mockTraceRepo)
	repo.On("GetTraces", mock.Anything, mock.MatchedBy(func(f repository.TraceFilter) bool {
		return f.AgentName == "test-agent"
	})).Return([]model.Trace{{TraceID: "1", AgentName: "test-agent", Timestamp: now}}, nil)

	h := handler.NewTraceHandler(repo)
	r := gin.New()
	rg := RouteRegistry{TraceHandler: h}
	RegisterRoutes(r, rg)

	req := httptest.NewRequest(http.MethodGet, "/api/traces?agent=test-agent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var res struct {
		Traces []model.Trace `json:"traces"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res.Traces, 1)
	assert.Equal(t, "1", res.Traces[0].TraceID)
	assert.Equal(t, "test-agent", res.Traces[0].AgentName)
}

func TestGetTraceByIDRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("not found", func(t *testing.T) {
		repo := new(mockTraceRepo)
		repo.On("GetByID", mock.Anything, "missing").Return((*model.Trace)(nil), errors.New("not found"))

		h := handler.NewTraceHandler(repo)
		r := gin.New()
		rg := RouteRegistry{TraceHandler: h}
		RegisterRoutes(r, rg)

		req := httptest.NewRequest(http.MethodGet, "/api/traces/missing", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error":"trace not found"}`, rec.Body.String())
		repo.AssertExpectations(t)
	})
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockTraceRepo)
	h := handler.NewTraceHandler(repo)
	registry := RouteRegistry{TraceHandler: h}

	r := gin.New()
	RegisterRoutes(r, registry)

	reqBody := model.Trace{TraceID: "xyz789", AgentName: "TestAgent"}
	body, _ := json.Marshal(reqBody)
	repo.On("InsertTrace", mock.Anything, mock.MatchedBy(func(t model.Trace) bool {
		return t.TraceID == reqBody.TraceID && t.AgentName == reqBody.AgentName
	})).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/traces", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"message":"trace saved"}`, rec.Body.String())
	repo.AssertExpectations(t)
}
