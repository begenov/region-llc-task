package service

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
)

type Users interface {
	SignUp(ctx context.Context, inp domain.UserRequest) (domain.User, error)
	SignIn(ctx context.Context, email, password string) (domain.Token, error)
	RefreshTokens(ctx context.Context, refreshToken string) (domain.Token, error)
}

type Todo interface {
	CreateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error)
}
