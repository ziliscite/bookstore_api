package models

import "time"

type Review struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Rating  int    `json:"rating" db:"rating"`
	Comment string `json:"comment" db:"comment"`

	UserID int `json:"user_id" db:"user_id"`
	BookID int `json:"book_id" db:"book_id"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
