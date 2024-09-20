package models

import "time"

// User represents the structure of the Users table in the database.
type User struct {
	ID        int64     `json:"id" db:"id"`                 // ID is the primary key of the user.
	Name      string    `json:"name" db:"name"`             // Name is the name of the user.
	Email     string    `json:"email" db:"email"`           // Email is the unique email address of the user.
	Password  string    `json:"password" db:"password"`     // Password is the hashed password of the user.
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`     // IsAdmin indicates if the user is an admin or not.
	CreatedAt time.Time `json:"created_at" db:"created_at"` // CreatedAt is the timestamp when the user was created.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // UpdatedAt is the timestamp when the user information was last updated.
}
