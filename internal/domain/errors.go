package domain

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("not found")
	ErrInternalServer     = errors.New("internal server")
)
