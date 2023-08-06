package mongo

import (
	"context"
	"fmt"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	collection := db.Collection(usersCollection)

	return &UserRepo{
		collection: collection,
	}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Errorf("r.collection.InsertOne(): %v", err)
		return domain.User{}, fmt.Errorf("r.collection.InsertOne(): %v", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
		return domain.User{}, fmt.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
	}

	user.ID = insertedID

	return user, nil
}

func (r *UserRepo) GetUserByUserName(ctx context.Context, username string) (domain.User, error) {
	return domain.User{}, nil
}

func GetUserByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	return domain.User{}, nil
}

func SetSession(ctx context.Context, id primitive.ObjectID, session domain.Session) (domain.User, error) {
	return domain.User{}, nil
}
