package router

import (
	"github.com/gin-gonic/gin"

	"github.com/zkropotkine/agent-trace/internal/handler"
)

type RouteRegistry struct {
	TraceHandler handler.TraceHandler
}

func SetupRouter(deps RouteRegistry) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, deps)

	return router
}
