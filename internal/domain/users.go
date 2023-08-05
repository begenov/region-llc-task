package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id"`
	UserName string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Session  Session            `json:"-"`
	CreateAt time.Time          `json:"create_at"`
	UpdateAt time.Time          `json:"update_at"`
}
