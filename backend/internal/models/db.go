package models

import "time"

type Book struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	ISBN            string    `json:"isbn" db:"isbn"`
	PublicationYear int       `json:"publication_year" db:"publication_year"`
	Available       bool      `json:"available" db:"available"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBookRequest struct {
	Title           string `json:"title" validate:"required"`
	Author          string `json:"author" validate:"required"`
	ISBN            string `json:"isbn" validate:"required"`
	PublicationYear int    `json:"publication_year"`
}

type UpdateBookRequest struct {
	Title           *string `json:"title,omitempty"`
	Author          *string `json:"author,omitempty"`
	ISBN            *string `json:"isbn,omitempty"`
	PublicationYear *int    `json:"publication_year,omitempty"`
	Available       *bool   `json:"available,omitempty"`
}
