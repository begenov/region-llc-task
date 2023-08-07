package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	TodoID   string             `json:"id" bson:"id"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title    string             `json:"title" bson:"title"`
	ActiveAt string             `json:"activeAt" bson:"activeAt"`
	Actor    string             `json:"actor" bson:"actor"`
	Status   string             `json:"status" bson:"status"`
}

type TodoRequest struct {
	Title    string `json:"title" binding:"required"`
	ActiveAt string `json:"activeAt" binding:"required"`
}

type TodoURI struct {
	ID string `uri:"id" binding:"required"`
}
