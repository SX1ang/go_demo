package redis

import (
	"context"
	"demoProject/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}

	return client, nil
}
