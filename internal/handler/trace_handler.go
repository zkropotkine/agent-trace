package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zkropotkine/agent-trace/pkg/logger"
	"github.com/zkropotkine/agent-trace/pkg/tokens"

	"github.com/zkropotkine/agent-trace/internal/model"
	"github.com/zkropotkine/agent-trace/internal/repository"
)

type traceHandler struct {
	repo repository.TraceRepository
}

func NewTraceHandler(repo repository.TraceRepository) TraceHandler {
	return &traceHandler{repo: repo}
}

func (h *traceHandler) PostTrace(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var trace model.Trace
	if err := c.ShouldBindJSON(&trace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trace payload"})
		return
	}
	trace.Timestamp = time.Now()

	info, err := tokens.Analyze(trace.InputPrompt, trace.OutputPrompt, trace.Model)

	if err != nil {
		log.Error("Error analyzing tokens:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze tokens"})
		return
	}

	trace.TokenUsage = model.TokenUsage{
		Input:         info.InputTokens,
		Output:        info.OutputTokens,
		Total:         info.TotalTokens,
		EstimatedCost: info.EstimatedCost,
	}

	err = h.repo.InsertTrace(c.Request.Context(), trace)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save trace"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "trace saved"})
}

func (h *traceHandler) GetTraces(c *gin.Context) {
	agent := c.Query("agent")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	var from, to *time.Time
	if fromStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &parsed
		}
	}
	if toStr != "" {
		parsed, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &parsed
		}
	}

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	filter := repository.TraceFilter{
		AgentName: agent,
		From:      from,
		To:        to,
		Limit:     limit,
		Offset:    offset,
	}

	traces, err := h.repo.GetTraces(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch traces"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"traces": traces})
}

func (h *traceHandler) GetTraceByID(c *gin.Context) {
	id := c.Param("id")

	trace, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trace not found"})
		return
	}

	c.JSON(http.StatusOK, trace)
}
