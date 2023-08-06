package v1

import (
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, Response{message})
}
