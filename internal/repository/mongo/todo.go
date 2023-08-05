package mongo

import "go.mongodb.org/mongo-driver/mongo"

type TodoRepo struct {
}

func NewTodoRepo(db *mongo.Database) *TodoRepo {
	return &TodoRepo{}
}
