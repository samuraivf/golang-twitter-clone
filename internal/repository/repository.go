package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user dto.CreateUserDto) (uint, error)
	EditProfile(user dto.EditUserDto) error
	AddImage(image []byte, userId uint) error
}

type Redis interface {
	SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
	GetUserRefreshTokens(ctx context.Context, pattern string) ([]string, error)
	DeleteUserRefreshTokens(ctx context.Context, keys []string) error
}

type Tweet interface {
	CreateTweet(tweetDto dto.CreateTweetDto, mentionedUsers []string) (uint, error)
	GetTweetById(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]*models.Tweet, error)
	UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error)
	DeleteTweet(tweetId uint) error
	LikeTweet(tweetId, userId uint) error
	UnlikeTweet(tweetId, userId uint) error
}

type Comment interface {
	CreateComment(commentDto dto.CreateCommentDto) (uint, error)
	GetCommentById(id uint) (*models.Comment, error)
	UpdateComment(commentDto dto.UpdateCommentDto) (uint, error)
	DeleteComment(id uint) error
	LikeComment(commentId, userId uint) error
	UnlikeComment(commentId, userId uint) error
}

type Repository struct {
	User
	Redis
	Tweet
	Comment
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		User:  NewUserPostgres(db),
		Tweet: NewTweetPostgres(db),
		Redis: NewRedis(redis),
		Comment: NewCommentPostgres(db),
	}
}
