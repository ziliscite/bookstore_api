package services

import (
	"encoding/base64"
	"errors"
	"os"
)

type Service struct {
	AESKey []byte
}

func NewService() (*Service, error) {
	encodedKey := os.Getenv("AES_KEY")
	if encodedKey == "" {
		return nil, errors.New("AES_KEY environment variable not set")
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, errors.New("failed to decode base64 key")
	}

	return &Service{
		AESKey: key,
	}, nil
}
