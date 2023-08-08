package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s *Server) createTodo(ctx *gin.Context) {
	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	if err := validateTodo(inp.Title, inp.ActiveAt); err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error(), fmt.Sprintf("validateTodo(): %v", err))
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

func (s *Server) updateTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindUri(): %v", err))
		return
	}

	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error(), fmt.Sprintf("ctx.BindJSON(): %v", err))
		return
	}

	if err := validateTodo(inp.Title, inp.ActiveAt); err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error(), fmt.Sprintf("validateTodo(): %v", err))
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		newResponse(ctx, checkErrors(err), err.Error(), fmt.Sprintf("getUserID(): %v", err))
		return
	}

	todo := domain.Todo{
		TodoID:   uri.ID,
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

func validateTodo(title string, activeAt string) error {

	if title == "" {
		return domain.ErrInvalidTitle
	}

	if len(title) > 200 {
		return domain.ErrHeaderLength
	}

	_, err := time.Parse("2006-01-02", activeAt)
	if err != nil {
		return domain.ErrIncorrectDateFormat
	}

	return nil
}
