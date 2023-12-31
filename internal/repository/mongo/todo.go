package mongo

import (
	"context"
	"errors"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
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
		return domain.Todo{}, domain.ErrInternalServer
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
		return domain.Todo{}, domain.ErrInternalServer
	}

	todo.ID = id

	return todo, nil
}

func (r *TodoRepo) GetCountByTitle(ctx context.Context, title string, id primitive.ObjectID) (int64, error) {
	c, err := r.collection.CountDocuments(ctx, bson.M{"title": title, "user_id": id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, domain.ErrNotFound
	}

	return c, err
}

func (r *TodoRepo) GetTodoByID(ctx context.Context, id primitive.ObjectID) (domain.Todo, error) {
	var todo domain.Todo
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo); err != nil {
		logger.Errorf("r.collection.FindOne(): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Todo{}, domain.ErrNotFound
		}

		return domain.Todo{}, err
	}

	return todo, nil
}

func (r *TodoRepo) UpdateTodoID(ctx context.Context, todo domain.Todo) error {
	_, err := r.collection.UpdateByID(ctx, todo.ID, bson.M{"$set": bson.M{"_id": todo.ID}})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrNotFound
	}

	return err
}

func (r *TodoRepo) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	_, err := r.collection.UpdateByID(ctx, todo.ID, bson.M{"$set": bson.M{"title": todo.Title, "activeAt": todo.ActiveAt}})
	if err != nil {
		logger.Errorf("r.collection.UpdateByID(): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *TodoRepo) DeleteTodoByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Errorf("r.collection.DeleteOne(): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *TodoRepo) UpdateTodoDoneByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (domain.Todo, error) {
	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{"$set": bson.M{"status": domain.Done}}

	var updatedTodo domain.Todo
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedTodo)
	if err != nil {
		logger.Errorf("r.collection.FindOneAndUpdate(): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Todo{}, domain.ErrNotFound
		}

		return domain.Todo{}, err
	}

	return updatedTodo, nil
}

func (r *TodoRepo) GetTodoByStatus(ctx context.Context, status string, userID primitive.ObjectID) ([]domain.Todo, error) {
	var todos []domain.Todo
	cur, err := r.collection.Find(ctx, bson.M{"status": status, "user_id": userID})
	if err != nil {
		logger.Errorf("r.collection.Find(): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var todo domain.Todo
		if err := cur.Decode(&todo); err != nil {
			logger.Errorf("cur.Decode(): %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := cur.Err(); err != nil {
		logger.Errorf("cur.Err(): %v", err)
		return nil, err
	}

	return todos, nil
}
