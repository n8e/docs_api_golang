package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentSchema struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	OwnerId      primitive.ObjectID `json:"ownerId"`
	Title        string             `json:"title,omitempty" validate:"required"`
	Content      string             `json:"content,omitempty" validate:"required"`
	DateCreated  string             `json:"dateCreated,omitempty"`
	LastModified string             `json:"lastModified,omitempty"`
}
