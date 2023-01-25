package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentSchema struct {
	OwnerId      primitive.ObjectID `json:"ownerId" validate:"required"`
	Title        string             `json:"title,omitempty" validate:"required"`
	Content      string             `json:"content,omitempty" validate:"required"`
	DateCreated  string             `json:"dateCreated,omitempty" validate:"required"`
	LastModified string             `json:"lastModified,omitempty" validate:"required"`
}
