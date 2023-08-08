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
	Author   string             `json:"author" bson:"author"`
	Status   string             `json:"status" bson:"status"`
}

type TodoRequest struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
}

type TodoURI struct {
	ID string `uri:"id" binding:"required"`
}
