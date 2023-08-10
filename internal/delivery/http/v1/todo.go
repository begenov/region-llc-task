package v1

import (
	"fmt"
	"net/http"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary		User Create New Todo List
// @Security UserAuth
// @Tags			Todo List
// @Description	User Create New Todo List
// @Accept			json
// @Produce		json
// @Param			account	body		domain.TodoRequest	true	"Todo-List"
// @Success		200		{object}	domain.Todo
// @Failure		400		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/todo-list/todo [post]
func (s *Server) createTodo(ctx *gin.Context) {
	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("getUserID(): %v", err))
		return
	}

	todo := domain.Todo{
		UserID:   id,
		Title:    inp.Title,
		ActiveAt: inp.ActiveAt,
		Status:   domain.Active,
	}

	todo, err = s.todoService.CreateTodo(ctx, todo)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.todoService.CreateTodo(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// @Summary		User Update Todo List
// @Security UserAuth
// @Tags			Todo List
// @Description	User Update Todo List
// @Accept			json
// @Produce		json
// @Param			account	body		domain.TodoRequest	true	"Todo List update info"
// @Param			domain.TodoURI path string		true	"Todo List id"
// @Success		200		{object}	domain.Todo
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/todo-list/{id} [put]
func (s *Server) updateTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindUri(): %v", err))
		return
	}

	todoID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindUri(): %v", err))
		return
	}

	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("getUserID(): %v", err))
		return
	}

	todo := domain.Todo{
		ID:       todoID,
		UserID:   id,
		Title:    inp.Title,
		ActiveAt: inp.ActiveAt,
		Status:   domain.Active,
	}

	todo, err = s.todoService.UpdateTodo(ctx, todo)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.todoService.UpdateTodo(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// @Summary		User Delete Todo List
// @Security UserAuth
// @Tags			Todo List
// @Description	User Delete Todo List
// @Accept			json
// @Produce		json
// @Param			domain.TodoURI path string		true	"Todo List ID"
// @Success		200		{object}	Response
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/todo-list/{id} [delete]
func (s *Server) deleteTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindUri(): %v", err))
		return
	}

	err := s.todoService.DeleteTodoByID(ctx, uri.ID)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.todoService.DeleteTodoByID(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, Response{"Success Deleting Todo"})
}

// @Summary		User Update Todo List
// @Security UserAuth
// @Tags			Todo List
// @Description	User Update Todo List
// @Accept			json
// @Produce		json
// @Param			domain.TodoURI path string		true	"Todo List id"
// @Success		200		{object}	domain.Todo
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/todo-list/{id}/done [put]
func (s *Server) doneTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindUri(): %v", err))
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("getUserID(): %v", err))
		return
	}

	todo, err := s.todoService.UpdateTodoDoneByID(ctx, uri.ID, id)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.todoService.UpdateTodoDoneByID(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// @Summary		User Get Todo List By Status
// @Security UserAuth
// @Tags			Todo List
// @Description	User Get Todo List By Status
// @Accept			json
// @Produce		json
// @Param	status	query	string	false	"Task status (default: active)"
// @Success		201		{object}	[]domain.Todo
// @Failure		400		{object}	Response
// @Failure		404		{object}	Response
// @Failure		500		{object}	Response
// @Router			/users/todo-list/todo{id} [get]
func (s *Server) getTodos(ctx *gin.Context) {
	status := ctx.DefaultQuery("status", "active")

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error(), fmt.Sprintf("getUserID(): %v", err))
		return
	}

	tasks, err := s.todoService.GetTodosByStatus(ctx, status, id)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("s.todoService.GetTodosByStatus(): %v", err))
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
