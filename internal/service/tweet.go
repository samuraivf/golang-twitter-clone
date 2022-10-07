package service

import (
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
)

type TweetService struct {
	repo repository.Tweet
}

func NewTweetService(repo repository.Tweet) *TweetService {
	return &TweetService{repo}
}

func (s *TweetService) CreateTweet(tweetDto dto.CreateTweetDto) (uint, error) {
	return s.repo.CreateTweet(tweetDto)
}

func (s *TweetService) GetTweetById(id uint) (*models.Tweet, error) {
	return s.repo.GetTweetById(id)
}

func (s *TweetService) GetUserTweets(userId uint) ([]*models.Tweet, error) {
	return s.repo.GetUserTweets(userId)
}

func (s *TweetService) UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error) {
	return s.repo.UpdateTweet(tweetDto)
}

func (s *TweetService) DeleteTweet(tweetId uint) error {
	return s.repo.DeleteTweet(tweetId)
}