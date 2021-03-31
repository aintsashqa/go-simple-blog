package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	err := client.Ping(ctx).Err()
	return client, err
}
