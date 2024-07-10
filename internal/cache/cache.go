package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("cache not found")

type Cache[T any] interface {
	One(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
