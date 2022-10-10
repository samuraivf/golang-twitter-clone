package service

import (
	"context"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
)

type Authorization interface {
	GenerateAccessToken(username string, userId uint) (string, error)
	GenerateRefreshToken(username string, userId uint) (*RefreshTokenData, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	ValidateUser(username, password string) (*models.User, error)
	CreateUser(user dto.CreateUserDto) (uint, error)
	GetUserByUsername(username string) (*models.User, error)
	EditProfile(user dto.EditUserDto) error
	AddImage(image string, userId uint) error
}

type Redis interface {
	SetRefreshToken(ctx context.Context, key, refreshToken string) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
}

type Tweet interface {
	CreateTweet(tweetDto dto.CreateTweetDto) (uint, error)
	GetTweetById(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]*models.Tweet, error)
	UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error)
	DeleteTweet(tweetId uint) error
	LikeTweet(tweetId, userId uint) error
	UnlikeTweet(tweetId, userId uint) error
}

type Service struct {
	Authorization
	User
	Redis
	Tweet
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(),
		User:          NewUserService(repos.User),
		Redis:         NewRedisService(repos.Redis),
		Tweet: 		   NewTweetService(repos.Tweet),
	}
}
