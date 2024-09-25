package models

import "time"

type Sessions struct {
	ID           string    `json:"id" db:"id"`
	UserEmail    string    `json:"user_email" db:"user_email"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked" db:"is_revoked"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
