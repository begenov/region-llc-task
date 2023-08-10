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

// @Summary		Create a new User
// @Tags			User
// @Description	Create a new User with the input payload
// @Accept			json
// @Produce		json
// @Param			account	body		domain.UserRequest	true	"User"
// @Success		200		{object}	Response
// @Failure		400		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/sign-up [post]
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

// @Summary		Sign-in
// @Tags			User
// @Description	Sign-in
// @Accept			json
// @Produce		json
// @Param			account	body		domain.UserSignInRequest	true	"User"
// @Success		200		{object}	Response
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/sign-in [post]
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

// @Summary		Refresh Token
// @Tags			User
// @Description	Refresh Token
// @Accept			json
// @Produce		json
// @Param			account	body		domain.RefreshToken	true	"User"
// @Success		200		{object}	domain.Token
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/auth/refresh [post]
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
