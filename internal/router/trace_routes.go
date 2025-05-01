package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterTraceRoutes(rg *gin.RouterGroup, deps RouteRegistry) {
	rg.POST("/traces", deps.TraceHandler.PostTrace)
}
