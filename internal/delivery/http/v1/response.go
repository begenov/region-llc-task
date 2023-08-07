package v1

import (
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string, log string) {
	logger.Error(log)
	c.AbortWithStatusJSON(statusCode, Response{message})
}
