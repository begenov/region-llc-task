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
		users.POST("/sign-in", s.userSignIn)
		users.POST("/auth/refresh", s.userRefresh)
		authenticated := users.Group("/", s.userIdentity)
		{
			todo := authenticated.Group("/todo-list")
			{
				todo.POST("/todo", s.createTodo)
				todo.PUT("/todo/:id", s.updateTodo)
				todo.DELETE("/todo/:id", s.deleteTodo)
				todo.DELETE("/todo/:id/done", s.doneTodo)

			}
		}
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

func (s *Server) userSignIn(ctx *gin.Context) {
	var inp domain.UserSignInRequest
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Message: domain.ErrInvalidRequest.Error()})
		return
	}

	tokens, err := s.userService.SignIn(ctx, inp.Email, inp.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrInternalServer.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (s *Server) userRefresh(c *gin.Context) {
	var inp domain.RefreshToken
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := s.userService.RefreshTokens(c.Request.Context(), inp.RefreshToken)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
