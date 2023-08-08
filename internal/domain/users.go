package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Session  Session            `bson:"session,omitempty" json:"-"`
	CreateAt string             `bson:"create_at" json:"create_at"`
}

type UserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
