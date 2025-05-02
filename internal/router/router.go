package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zkropotkine/agent-trace/internal/middleware"
	"github.com/zkropotkine/agent-trace/pkg/logger"

	"github.com/zkropotkine/agent-trace/internal/handler"
)

type RouteRegistry struct {
	TraceHandler handler.TraceHandler
}

func SetupRouter(ctx context.Context, registry RouteRegistry) *gin.Engine {
	log := logger.FromContext(ctx)
	router := gin.Default()

	router.Use(
		gin.Recovery(), // catches panics and logs stack traces
		middleware.RequestLogger(log),
	)

	RegisterRoutes(router, registry)

	return router
}
