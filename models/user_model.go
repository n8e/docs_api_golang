package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoleType string

const (
	User          RoleType = "User"
	Administrator RoleType = "Administrator"
)

type UserSchema struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	UserName  string             `json:"username,omitempty" validate:"required"`
	FirstName string             `json:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" validate:"required,email"`
	Password  string             `json:"password,omitempty"`
	Role      RoleType           `json:"role" validate:"required,oneof=User Administrator"`
}
