package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, deps RouteRegistry) {
	api := router.Group("/api")
	{
		api.POST("/traces", deps.TraceHandler.PostTrace)
		// RegisterEvaluationRoutes(api, deps) ← future
	}
}
