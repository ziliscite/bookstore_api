package models

import "time"

// Book represents the structure of the Books table in the database.
type Book struct {
	ID    int    `json:"id" db:"id"`       // ID is the primary key of the book.
	Title string `json:"title" db:"title"` // Title is the unique title of the book.

	Slug       string `json:"slug" db:"slug"`               // Slug is a unique identifier for the book based off Title.
	CoverImage string `json:"cover_image" db:"cover_image"` // CoverImage holds the URL/path to the book's cover image.
	Synopsis   string `json:"synopsis" db:"synopsis"`       // Synopsis provides a description or summary of the book.

	Price float64 `json:"price" db:"price"` // Price is the cost of the book, with 2 decimal places.
	Stock int     `json:"stock" db:"stock"` // Stock represents how many copies of the book are available.

	CreatedAt time.Time  `json:"created_at" db:"created_at"` // CreatedAt holds the timestamp when the book was created.
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"` // UpdatedAt holds the timestamp when the book was last updated. Nullable.
}

// json tag is needed for marshalling
// db tag is to map the field to the database queries
