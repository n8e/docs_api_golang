package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentSchema struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	OwnerId      primitive.ObjectID `json:"ownerId"`
	Title        string             `json:"title,omitempty" validate:"required"`
	Content      string             `json:"content,omitempty" validate:"required"`
	DateCreated  time.Time          `json:"dateCreated,omitempty"`
	LastModified time.Time          `json:"lastModified,omitempty"`
}
