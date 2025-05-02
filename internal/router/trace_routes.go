package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, deps RouteRegistry) {
	api := router.Group("/api")
	{
		api.POST("/traces", deps.TraceHandler.PostTrace)
		api.GET("/traces", deps.TraceHandler.GetTraces)
		api.GET("/traces/:id", deps.TraceHandler.GetTraceByID)
		// RegisterEvaluationRoutes(api, deps) ← future
	}
}
