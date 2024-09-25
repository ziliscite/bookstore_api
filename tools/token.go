package tools

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

type ContextKey string

type CustomClaims struct {
	IsAdmin bool `json:"isAdmin,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, isAdmin bool, duration time.Duration) (*CustomClaims, string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, "", errors.New("error generating token")
	}

	claims := &CustomClaims{
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenID.String(),

			Subject: email,
			Issuer:  os.Getenv("ISSUER"),

			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}

	return claims, signedToken, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(key *jwt.Token) (interface{}, error) {
		if _, ok := key.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", key.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil || issuer != os.Getenv("ISSUER") {
		return nil, errors.New("invalid credentials")
	}

	if !token.Valid {
		return nil, errors.New("invalid credentials")
	}

	return claims, nil
}
