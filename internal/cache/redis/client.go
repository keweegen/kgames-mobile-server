package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheClient struct {
	keyPrefix string
	base      *redis.Client
}

func NewCacheClient(keyPrefix string, client *redis.Client) CacheClient {
	return CacheClient{
		keyPrefix: keyPrefix,
		base:      client,
	}
}

func (c CacheClient) key(v string) string {
	return fmt.Sprintf("%s:%s", c.keyPrefix, v)
}

func (c CacheClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.base.Get(ctx, c.key(key))
}

func (c CacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	return c.base.Set(ctx, c.key(key), value, ttl)
}

func (c CacheClient) Delete(ctx context.Context, key string) *redis.IntCmd {
	return c.base.Del(ctx, c.key(key))
}
