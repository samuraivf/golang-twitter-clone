package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

type Redis interface {
	SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
	GetUserRefreshTokens(ctx context.Context, pattern string) ([]string, error)
	DeleteUserRefreshTokens(ctx context.Context, keys []string) error
}

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

func (r *RedisRepository) GetUserRefreshTokens(ctx context.Context, pattern string) ([]string, error) {
	return r.redis.Keys(ctx, pattern).Result()
}

func (r *RedisRepository) DeleteUserRefreshTokens(ctx context.Context, keys []string) error {
	return r.redis.Del(ctx, keys...).Err()
}

func (r *RedisRepository) SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error {
	userTokens, err := r.GetUserRefreshTokens(ctx, fmt.Sprintf("*%s*", strings.Split(key, ":")[0]))
	if err != nil {
		return err
	}

	if len(userTokens) > 5 {
		err := r.DeleteUserRefreshTokens(ctx, userTokens)
		if err != nil {
			return err
		}
	}

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
