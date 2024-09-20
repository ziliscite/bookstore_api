package handlers

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Handler struct {
	Cache   *redis.Client
	Timeout time.Duration
}

func NewHandler(cache *redis.Client) *Handler {
	return &Handler{
		Cache:   cache,
		Timeout: 3600 * time.Second,
	}
}
