package service

import (
	"context"

	"github.com/samuraivf/twitter-clone/internal/repository"
)

type RedisService struct {
	repo repository.Redis
}

func NewRedisService(repo repository.Redis) *RedisService {
	return &RedisService{repo}
}

func (s *RedisService) SetRefreshToken(ctx context.Context, key, refreshToken string) error {
	return s.repo.SetRefreshToken(ctx, key, refreshToken, RefreshTokenTTL)
}

func (s *RedisService) GetRefreshToken(ctx context.Context, key string) (string, error) {
	return s.repo.GetRefreshToken(ctx, key)
}

func (s *RedisService) DeleteRefreshToken(ctx context.Context, key string) error {
	return s.repo.DeleteRefreshToken(ctx, key)
}
