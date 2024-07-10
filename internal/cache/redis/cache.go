package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/keweegen/tic-toe/internal/cache"
	"time"
)

type Cache[T any] struct {
	keyPrefix string
	client    CacheClient
}

var _ cache.Cache[any] = (*Cache[any])(nil)

func NewCache[T any](keyPrefix string, client CacheClient) Cache[T] {
	return Cache[T]{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

func (c Cache[T]) key(v string) string {
	return fmt.Sprintf("%s:%s", c.keyPrefix, v)
}

func (c Cache[T]) One(ctx context.Context, key string) (T, error) {
	var data T

	str, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return data, cache.ErrNotFound
		}
		return data, err
	}

	if err = json.Unmarshal([]byte(str), &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c Cache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, string(data), ttl).Err()
}

func (c Cache[T]) Delete(ctx context.Context, key string) error {
	return c.client.Delete(ctx, key).Err()
}
