package mongo

import (
	"context"
	"fmt"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepo struct {
	collection *mongo.Collection
}

func NewTodoRepo(db *mongo.Database) *TodoRepo {
	return &TodoRepo{db.Collection(todoCollection)}
}

func (r *TodoRepo) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	result, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		logger.Errorf("r.collection.InsertOne(): %v", err)
		return domain.Todo{}, fmt.Errorf("r.collection.InsertOne(): %v", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
		return domain.Todo{}, fmt.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
	}

	todo.ID = id

	return todo, nil
}

func (r *TodoRepo) GetTodoByID(ctx context.Context, id primitive.ObjectID) (domain.Todo, error) {
	return domain.Todo{}, nil
}

func (r *TodoRepo) UpdateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	return domain.Todo{}, nil
}
