package v1

import (
	"net/http"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) createTodo(ctx *gin.Context) {
	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		logger.Errorf("ctx.BindJSON(): %v", err)
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	if err := validateTodo(inp.Title, inp.ActiveAt); err != nil {
		logger.Errorf("validateTodo(): %v", err)
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		logger.Errorf("getUserID(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
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
		logger.Errorf("s.todoService.CreateTodo(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (s *Server) updateTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		logger.Errorf("ctx.BindUri(): %v", err)
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		logger.Errorf("ctx.BindJSON(): %v")
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	if err := validateTodo(inp.Title, inp.ActiveAt); err != nil {
		logger.Errorf("validateTodo(): %v", err)
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
		logger.Errorf("getUserID(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
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
		logger.Errorf("s.todoService.UpdateTodo(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (s *Server) deleteTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		logger.Errorf("ctx.BindUri(): %v", err)
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	err := s.todoService.DeleteTodoByID(ctx, uri.ID)
	if err != nil {
		logger.Errorf("primitive.ObjectIDFromHex(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Deleted todo")
}

func (s *Server) doneTodo(ctx *gin.Context) {
	var uri domain.TodoURI
	if err := ctx.BindUri(&uri); err != nil {
		logger.Errorf("ctx.BindUri(): %v", err)
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	todo, err := s.todoService.UpdateTodoDoneByID(ctx, uri.ID)
	if err != nil {
		logger.Errorf("primitive.ObjectIDFromHex(): %v", err)
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (s *Server) getTodos(ctx *gin.Context) {
	status := ctx.DefaultQuery("status", "active")

	// Call your service method to get tasks based on the status
	tasks, err := s.todoService.GetTodosByStatus(ctx, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)

}

func validateTodo(title string, activeAt string) error {
	if len(title) > 200 {
		return domain.ErrHeaderLength
	}

	_, err := time.Parse("2006-01-02", activeAt)
	if err != nil {
		return domain.ErrIncorrectDateFormat
	}

	return nil
}
