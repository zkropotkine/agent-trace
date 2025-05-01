package handler

import "github.com/gin-gonic/gin"

type TraceHandler interface {
	PostTrace(c *gin.Context)
}
