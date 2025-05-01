package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zkropotkine/agent-trace/internal/model"
	"github.com/zkropotkine/agent-trace/internal/repository"
)

type traceHandlerImpl struct {
	repo repository.TraceRepository
}

func NewTraceHandler(repo repository.TraceRepository) TraceHandler {
	return &traceHandlerImpl{repo: repo}
}

func (h *traceHandlerImpl) PostTrace(c *gin.Context) {
	var trace model.Trace
	if err := c.ShouldBindJSON(&trace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trace payload"})
		return
	}
	trace.Timestamp = time.Now()

	err := h.repo.InsertTrace(c.Request.Context(), trace)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save trace"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "trace saved"})
}
