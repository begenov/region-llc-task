package v1

import (
	"fmt"
	"net/http"

	"github.com/begenov/region-llc-task/internal/domain"
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
				todo.PUT("/todo/:id/done", s.doneTodo)
				todo.GET("/todo", s.getTodos)
			}
		}
	}
}

func (s *Server) userSignUp(ctx *gin.Context) {
	var req domain.UserRequest

	if err := ctx.BindJSON(&req); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	user, err := s.userService.SignUp(ctx, req)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) userSignIn(ctx *gin.Context) {
	var inp domain.UserSignInRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	tokens, err := s.userService.SignIn(ctx, inp.Email, inp.Password)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.userService.SignIn(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (s *Server) userRefresh(ctx *gin.Context) {
	var inp domain.RefreshToken
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	res, err := s.userService.RefreshTokens(ctx, inp.RefreshToken)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.userService.RefreshTokens(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
