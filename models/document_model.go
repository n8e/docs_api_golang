package models

type DocumentSchema struct {
	OwnerId      string `json:"id,omitempty" validate:"required"`
	Title        string `json:"title,omitempty" validate:"required"`
	Content      string `json:"content,omitempty" validate:"required"`
	DateCreated  string `json:"dateCreated,omitempty" validate:"required"`
	LastModified string `json:"lastModified,omitempty" validate:"required"`
}
