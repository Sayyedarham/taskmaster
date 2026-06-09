package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client}
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *CacheRepository) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CacheRepository) Delete(ctx context.Context, pattern string) error {
	// Simple exact delete for now; production uses SCAN + DEL
	return r.client.Del(ctx, pattern).Err()
}

func (r *CacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *CacheRepository) Expire(ctx context.Context, key string, seconds int) error {
	return r.client.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}
