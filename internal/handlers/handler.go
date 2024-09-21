package handlers

import (
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Cache *redis.Client
}

func NewHandler(cache *redis.Client) *Handler {
	return &Handler{
		Cache: cache,
	}
}
