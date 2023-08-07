package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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
		return domain.User{}, domain.ErrInternalServer
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Errorf("result.InsertedID.(primitive.ObjectID): %v", ok)
		return domain.User{}, domain.ErrInternalServer
	}

	user.ID = id

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		logger.Errorf("r.collection.FindOne(email): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, domain.ErrInternalServer
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	var user domain.User

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		logger.Errorf("r.collection.FindOne(email): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrNotFound
		}

		logger.Errorf("r.collection.FindOne(email): %v", err)
		return domain.User{}, domain.ErrInternalServer
	}

	return user, nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	if err := r.collection.FindOne(ctx, bson.M{
		"session.refresh_token": refreshToken,
		"session.expiration_at": bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		logger.Errorf("r.collection.FindOne(email): %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, domain.ErrInternalServer
	}

	return user, nil
}

func (r *UserRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrNotFound
	}

	return err
}
