package service

import (
	"strings"

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
	var mentionedUsers []string

	for _, letter := range strings.Split(tweetDto.Text, " ") {
		if letter[0] == byte('@') {
			mentionedUsers = append(mentionedUsers, letter[1:])
		}
	}

	return s.repo.CreateTweet(tweetDto, mentionedUsers)
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

func (s *TweetService) LikeTweet(tweetId, userId uint) error {
	return s.repo.LikeTweet(tweetId, userId)
}

func (s *TweetService) UnlikeTweet(tweetId, userId uint) error {
	return s.repo.UnlikeTweet(tweetId, userId)
}
