package repository

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type Todo interface {
}
