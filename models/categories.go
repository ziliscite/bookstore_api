package models

type Category struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type BookCategory struct {
	CategoryID int64 `json:"category_id" db:"category_id"`
	BookID     int64 `json:"book_id" db:"book_id"`
}
