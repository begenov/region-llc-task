package v1

import (
	"net/http"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) initLoadRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", s.userSignUp)
	}
}

func (s *Server) userSignUp(ctx *gin.Context) {
	var req domain.UserRequest

	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("ctx.BindJSON(): %v", err)
		ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	user, err := s.userService.SignUp(ctx, req)
	if err != nil {
		logger.Errorf("ctx.BindJSON(): %v", err)
		ctx.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
