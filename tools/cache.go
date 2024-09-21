package tools

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// GetOrSetCacheById takes a key and value type, returns an encoded value
// I don't want to do generics. This stops here.
func getOrSetCacheById(c *redis.Client, ctx context.Context, id int, v interface{}, fn func(context.Context, int) (interface{}, error)) ([]byte, error) {
	k := ctx.Value("cachedKey").(string)

	// Returns cache if present
	cachedJSON, err := c.Get(ctx, k).Result()
	if !errors.Is(err, redis.Nil) {
		err = json.Unmarshal([]byte(cachedJSON), &v)
		if err != nil {
			return nil, errors.New("failed to unmarshal cache")
		}

		encodedCache, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("failed to process data")
		}

		return encodedCache, nil
	}

	// Getting data from a service callback
	var dat interface{}
	dat, err = fn(ctx, id)
	if err != nil {
		return nil, errors.New("failed to load data from the database")
	}

	// Encode the data before adding to the cache
	encodedDat, err := json.Marshal(dat)
	if err != nil {
		return nil, errors.New("failed to process data")
	}

	err = c.SetEx(ctx, k, encodedDat, time.Second*3600).Err()
	if err != nil {
		return nil, errors.New("failed to cache data")
	}

	// Return the encoded data
	return encodedDat, nil
}
