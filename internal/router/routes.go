package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, deps RouteRegistry) {
	api := router.Group("/api")
	{
		RegisterTraceRoutes(api, deps)
		// RegisterEvaluationRoutes(api, deps) ← future
	}
}
