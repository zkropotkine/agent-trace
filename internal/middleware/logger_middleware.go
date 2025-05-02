package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zkropotkine/agent-trace/pkg/logger"
)

// RequestLogger injects a logger into the context with request-scoped fields.
func RequestLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		baseLog := log.WithFields(map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"request_id": requestID,
		})

		// Set logger into context
		ctx := logger.WithLogger(c.Request.Context(), baseLog)
		c.Request = c.Request.WithContext(ctx)

		// Set request ID as header for response
		c.Writer.Header().Set("X-Request-Id", requestID)

		c.Next()
	}
}
