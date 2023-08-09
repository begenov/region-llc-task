package v1

import (
	"net/http"

	"github.com/begenov/region-llc-task/internal/domain"
)

func checkErrors(err error) int {
	switch err {
	case domain.ErrInvalidRequest, domain.ErrEmailAlreadyExists, domain.ErrIncorrectDateFormat, domain.ErrHeaderLength,
		domain.ErrTitleAlreadyExists, domain.ErrInvalidTitle, domain.ErrTodoInvalidId:
		return http.StatusBadRequest
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
