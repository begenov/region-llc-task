package domain

import "errors"

var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrNotFound            = errors.New("not found")
	ErrInternalServer      = errors.New("internal server")
	ErrIncorrectDateFormat = errors.New("incorrect date format")
	ErrHeaderLength        = errors.New("header length exceeds 200 characters")
	ErrTitleAlreadyExists  = errors.New("title already exists")
)
