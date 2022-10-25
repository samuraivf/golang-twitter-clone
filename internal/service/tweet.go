package service

import (
	"regexp"
	"strings"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
)

type TweetService struct {
	repo repository.Tweet
	tag  repository.Tag
}

func NewTweetService(repo repository.Tweet, tag repository.Tag) *TweetService {
	return &TweetService{repo, tag}
}

func (s *TweetService) CreateTweet(tweetDto dto.CreateTweetDto) (uint, error) {
	var mentionedUsers []string

	for _, letter := range strings.Split(tweetDto.Text, " ") {
		if letter[0] == byte('@') {
			mentionedUsers = append(mentionedUsers, letter[1:])
		}
	}

	tweetId, err := s.repo.CreateTweet(tweetDto, mentionedUsers)

	if err != nil {
		return 0, err
	}

	for _, tag := range tweetDto.Tags {
		r, err := regexp.Compile("^[a-z0-9]+$")

		if err != nil {
			return 0, err
		}

		if r.Match([]byte(tag)) {
			s.tag.CreateTag(tag)

			err = s.tag.AddTweet(tag, tweetId)
			if err != nil {
				return 0, err
			}
		}
	}

	if err != nil {
		return 0, err
	}

	return tweetId, nil
}

func (s *TweetService) GetTweetById(id uint) (*models.Tweet, error) {
	return s.repo.GetTweetById(id)
}

func (s *TweetService) GetUserTweets(userId uint) ([]models.Tweet, error) {
	return s.repo.GetUserTweets(userId)
}

func (s *TweetService) UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error) {
	var mentionedUsers []string

	for _, letter := range strings.Split(tweetDto.Text, " ") {
		if letter[0] == byte('@') {
			mentionedUsers = append(mentionedUsers, letter[1:])
		}
	}

	tweetId, err := s.repo.UpdateTweet(tweetDto, mentionedUsers)

	if err != nil {
		return 0, err
	}

	for _, tag := range tweetDto.Tags {
		r, err := regexp.Compile("^[a-z0-9]+$")

		if err != nil {
			return 0, err
		}

		if r.Match([]byte(tag)) {
			s.tag.CreateTag(tag)

			err = s.tag.AddTweet(tag, tweetId)
			if err != nil {
				return 0, err
			}
		}
	}

	if err != nil {
		return 0, err
	}

	return tweetId, nil
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
