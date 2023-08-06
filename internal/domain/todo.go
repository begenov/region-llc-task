package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID       primitive.ObjectID `json:"id"`
	UserID   primitive.ObjectID `json:"user_id"`
	Title    string             `json:"title"`
	ActiveAt time.Time          `json:"active_at"`
	Actor    string             `json:"actor"`
	Status   string             `json:"status"`
}

type TodoRequest struct {
	Title    string    `json:"title" binding:"required"`
	ActiveAt time.Time `json:"active_at" binding:"required"`
}
