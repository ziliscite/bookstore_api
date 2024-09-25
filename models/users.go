package models

import "time"

// User represents the structure of the Users table in the database.
type User struct {
	ID        int64      `json:"id" db:"id"`                 // ID is the primary key of the user.
	Name      string     `json:"name" db:"name"`             // Name is the name of the user.
	Email     string     `json:"email" db:"email"`           // Email is the unique email address of the user.
	Password  string     `json:"password" db:"password"`     // Password is the hashed password of the user.
	IsAdmin   bool       `json:"is_admin" db:"is_admin"`     // IsAdmin indicates if the user is an admin or not.
	CreatedAt time.Time  `json:"created_at" db:"created_at"` // CreatedAt is the timestamp when the user was created.
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"` // UpdatedAt is the timestamp when the user information was last updated.
}

// UserRegister represents the data user require providing when registering a new account
type UserRegister struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserResponse struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Email   string `json:"email" db:"email"`
	IsAdmin bool   `json:"is_admin" db:"is_admin"`
}

// UserLogin represents the data user require providing when logging in
type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserLoginResponse struct {
	User                  UserResponse `json:"user"`
	SessionsId            string       `json:"session_id" db:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
}

// UserUpdateData Only name for now
type UserUpdateData struct {
	Name string `json:"name" db:"name"`
}
