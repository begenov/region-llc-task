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
