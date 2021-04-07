package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisProvider struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisProvider(ctx context.Context, cfg Config, exp time.Duration) (*RedisProvider, error) {
	client, err := NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &RedisProvider{client: client, exp: exp}, nil
}

func (p *RedisProvider) Set(ctx context.Context, key string, value []byte) error {
	return p.client.Set(ctx, key, value, p.exp).Err()
}

func (p *RedisProvider) Get(ctx context.Context, key string) ([]byte, error) {
	return p.client.Get(ctx, key).Bytes()
}

func (p *RedisProvider) Delete(ctx context.Context, key string) error {
	return p.client.Del(ctx, key).Err()
}
