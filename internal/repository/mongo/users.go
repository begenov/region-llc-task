package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection(usersCollection),
	}
}

func (r *UserRepo) Create(ctx context.Context) {

}
