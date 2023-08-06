package v1

import (
	"net/http"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s *Server) createTodo(ctx *gin.Context) {
	var inp domain.TodoRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	id, err := getUserID(ctx, userCtx)
	if err != nil {
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
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (s *Server) updateTodo(ctx *gin.Context) {
}
