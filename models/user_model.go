package models

type Role string

const (
	User          Role = "User"
	Administrator      = "Administrator"
)

type UserSchema struct {
	UserName  string `json:"username,omitempty" validate:"required"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty" validate:"required"`
	Password  string `json:"password,omitempty"`
	Role      *Role  `json:"role,omitempty"`
}
