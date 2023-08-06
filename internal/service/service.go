package service

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
)

type Users interface {
	SignUp(ctx context.Context, inp domain.UserRequest) (domain.User, error)
}

type Todo interface {
}
