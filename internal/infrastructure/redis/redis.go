package redis

import (
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Cache struct {
	cache *redis.Client
}

func New() *Cache {
	return &Cache{cache: redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
	})}
}

func (c *Cache) Close() error {
	return c.cache.Close()
}

func (c *Cache) GetCache() *redis.Client {
	return c.cache
}

// docker run --name redis-server -p 6379:6379 -d redis
