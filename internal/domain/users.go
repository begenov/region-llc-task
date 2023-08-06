package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Session  Session            `bson:"session,omitempty" json:"-"`
	CreateAt time.Time          `bson:"create_at" json:"create_at"`
	UpdateAt time.Time          `bson:"update_at" json:"update_at"`
}

type UserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
