package repository

import (
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

type TweetPostgres struct {
	db *gorm.DB
}

func NewTweetPostgres(db *gorm.DB) *TweetPostgres {
	return &TweetPostgres{db}
}

func (r *TweetPostgres) CreateTweet(tweetDto dto.CreateTweetDto) (uint, error) {
	tweet := models.Tweet{
		Text: tweetDto.Text,
		UserID: tweetDto.UserID,
	}

	result := r.db.Create(&tweet)
	if result.Error != nil {
		return 0, result.Error
	}

	return tweet.ID, nil
}

func (r *TweetPostgres) UpdateTweet(tweetDto dto.UpdateTweetDto) (uint, error) {
	tweet, err := r.GetTweetById(tweetDto.TweetID)

	if err != nil {
		return  0, err
	}

	tweet.Text = tweetDto.Text

	if err := r.db.Save(&tweet).Error; err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

func (r *TweetPostgres) DeleteTweet(tweetId uint) error {
	return r.db.Delete(&models.Tweet{}, tweetId).Error
}

func (r *TweetPostgres) GetTweetById(id uint) (*models.Tweet, error) {
	var tweet *models.Tweet

	result := r.db.First(&tweet, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return tweet, nil
}

func (r *TweetPostgres) GetUserTweets(userId uint) ([]*models.Tweet, error) {
	var tweets []*models.Tweet

	result := r.db.Where("user_id = ?", userId).Find(&tweets)
	if result.Error != nil {
		return nil, result.Error
	}

	return tweets, nil
}