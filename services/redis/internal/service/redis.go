package service

import (
	"context"

	"redis/internal/repo"
	
	pb "redis/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type RedisService struct {
	pb.UnimplementedRedisServer
	repo repository.Redis
}

func NewRedisService(repo repository.Redis) *RedisService {
	return &RedisService{repo: repo}
}

func (s *RedisService) SetRefreshToken(ctx context.Context, params *pb.SetRefreshTokenParams) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.SetRefreshToken(ctx, params.Key, params.RefreshToken, params.RefreshTokenTTL.AsDuration())
}

func (s *RedisService) GetRefreshToken(ctx context.Context, key *pb.Key) (*pb.RefreshToken, error) {
	token, err := s.repo.GetRefreshToken(ctx, key.Key)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshToken{
		RefreshToken: token,
	}, nil
}

func (s *RedisService) DeleteRefreshToken(ctx context.Context, key *pb.Key) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.DeleteRefreshToken(ctx, key.Key)
}
