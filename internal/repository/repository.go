package repository

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}

type Todo interface {
	Create(ctx context.Context, todo domain.Todo) (domain.Todo, error)
	GetTodoByID(ctx context.Context, id primitive.ObjectID) (domain.Todo, error)
	UpdateTodo(ctx context.Context, todo domain.Todo) error
	GetCountByTitle(ctx context.Context, title string, id primitive.ObjectID) (int64, error)
	DeleteTodoByID(ctx context.Context, id primitive.ObjectID) error
}
