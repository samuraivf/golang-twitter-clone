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
	Subscribe(subscriberId, userId uint) error
	Unsubscribe(subscriberId, userId uint) error
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
	GetUserTweets(userId uint) ([]models.Tweet, error)
	UpdateTweet(tweetDto dto.UpdateTweetDto, mentionedUsers []string) (uint, error)
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

type Tag interface {
	CreateTag(name string) (uint, error)
	GetTagByName(name string) (*models.Tag, error)
	AddTweet(tagName string, tweetId uint) error
	GetTop100Tags() ([]*models.Tag, error)
	GetTagByIdWithTweets(id uint) (*models.Tag, error)
}

type Message interface {
	AddMessage(message *models.Message) error
	GetUserMessages(userId uint) ([]*models.Message, error)
}

type Repository struct {
	User
	Redis
	Tweet
	Comment
	Tag
	Message
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	message := NewMessagePostgres(db)

	return &Repository{
		User:    NewUserPostgres(db, message),
		Tweet:   NewTweetPostgres(db, message),
		Redis:   NewRedis(redis),
		Comment: NewCommentPostgres(db, message),
		Tag:     NewTagPostgres(db),
		Message: message,
	}
}
