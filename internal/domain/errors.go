package domain

import "errors"

var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrNotFound              = errors.New("not found")
	ErrInternalServer        = errors.New("internal server")
	ErrIncorrectDateFormat   = errors.New("incorrect date format")
	ErrHeaderLength          = errors.New("header length exceeds 200 characters")
	ErrTitleAlreadyExists    = errors.New("title already exists")
	ErrIncorrectEmailAddress = errors.New("incorrect email address")
	ErrIncorrectUserName     = errors.New("incorrect username")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrInvalidTitle          = errors.New("invalid empty title")
	ErrInvalidAuthHeader     = errors.New("invalid auth header")
	ErrTodoInvalidId         = errors.New("invalid todo id")
	ErrTodoActiveAtData      = errors.New("active date has already passed")
)
