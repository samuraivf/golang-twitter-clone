package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Host string
	Port string
}

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisClient(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	logrus.Printf("redis started: %s", pong)

	return client, nil
}

func NewRedis(redis *redis.Client) *RedisRepository {
	return &RedisRepository{redis}
}

func (r *RedisRepository) SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error {
	return r.redis.Set(ctx, key, refreshToken, TTL).Err()
}

func (r *RedisRepository) GetRefreshToken(ctx context.Context, key string) (string, error) {
	refreshToken, err := r.redis.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (r *RedisRepository) DeleteRefreshToken(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}
