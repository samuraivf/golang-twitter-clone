package repository

import (
	"errors"
	"fmt"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errUserAlreadyLikedATweet = errors.New("user have already liked a tweet")
)

type TweetPostgres struct {
	db *gorm.DB
	message Message
}

func NewTweetPostgres(db *gorm.DB, message Message) *TweetPostgres {
	return &TweetPostgres{db, message}
}

func (r *TweetPostgres) CreateTweet(tweetDto dto.CreateTweetDto, mentionedUsers []string) (uint, error) {
	tweet := models.Tweet{
		Text:   tweetDto.Text,
		UserID: tweetDto.UserID,
	}

	result := r.db.Create(&tweet)

	if result.Error != nil {
		return 0, result.Error
	}

	if err := r.addMentions(&tweet, mentionedUsers); err != nil {
		return 0, err
	}

	if err := r.notifySubscribers(tweet.UserID, tweet.ID); err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

func (r *TweetPostgres) notifySubscribers(authorId, tweetId uint) error {
	var tweetAuthor *models.User

	if err := r.db.Where("id = ?", authorId).Preload("Subscribers").First(&tweetAuthor).Error; err != nil {
		return err
	}
	
	for _, subscriber := range tweetAuthor.Subscribers {
		message := models.Message{
			Text: fmt.Sprintf("@%s posted a new tweet", tweetAuthor.Username),
			UserID: subscriber.ID,
			AuthorID: tweetAuthor.ID,
			TweetID: tweetId,
		}

		if err := r.message.AddMessage(&message); err != nil {
			return err
		}

	}

	return nil
}

func (r *TweetPostgres) UpdateTweet(tweetDto dto.UpdateTweetDto, mentionedUsers []string) (uint, error) {
	tweet, err := r.GetTweetById(tweetDto.TweetID)

	if err != nil {
		return 0, err
	}

	tweet.Text = tweetDto.Text

	if err := r.addMentions(tweet, mentionedUsers); err != nil {
		return 0, err
	}

	if err := r.db.Save(&tweet).Error; err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

func (r *TweetPostgres) addMentions(tweet *models.Tweet, mentionedUsers []string) error {
	var tweetAuthor *models.User

	if err := r.db.Select("username").First(&tweetAuthor, tweet.UserID).Error; err != nil {
		return err
	}

	for _, user := range mentionedUsers {
		var userFromDB models.User

		if err := r.db.Where("username = ?", user).First(&userFromDB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			} else {
				return err
			}
		}
		tweet.MentionedUsers = append(tweet.MentionedUsers, &userFromDB)

		message := models.Message{
			Text: fmt.Sprintf("You was metioned by @%s", tweetAuthor.Username),
			UserID: userFromDB.ID,
			AuthorID: tweet.UserID,
			TweetID: tweet.ID,
		}

		if err := r.message.AddMessage(&message); err != nil {
			return err
		}
	}

	return nil
}

func (r *TweetPostgres) DeleteTweet(tweetId uint) error {
	err := r.db.Exec("DELETE FROM tag_tweets WHERE tweet_id IN (SELECT id FROM tweets WHERE tweets.id = ?);", tweetId).Error
	if err != nil {
		return err
	}

	return r.db.Exec("DELETE FROM tweets WHERE tweets.id = ?", tweetId).Error
}

func (r *TweetPostgres) GetTweetById(id uint) (*models.Tweet, error) {
	var tweet *models.Tweet

	result := r.db.Where("id = ?", id).Preload(clause.Associations).First(&tweet)
	if result.Error != nil {
		return nil, result.Error
	}

	return tweet, nil
}

func (r *TweetPostgres) GetUserTweets(userId uint) ([]models.Tweet, error) {
	var tweets []models.Tweet

	result := r.db.Where("user_id = ?", userId).Find(&tweets)
	if result.Error != nil {
		return nil, result.Error
	}

	return tweets, nil
}

func (r *TweetPostgres) LikeTweet(tweetId, userId uint) error {
	var tweet models.Tweet

	if err := r.db.Where("id = ?", tweetId).Preload("Likes").First(&tweet).Error; err != nil {
		return err
	}

	for _, user := range tweet.Likes {
		if user.ID == userId {
			return errUserAlreadyLikedATweet
		}
	}

	var user models.User

	if err := r.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := r.db.Model(&tweet).Association("Likes").Append(&user); err != nil {
		return err
	}

	message := models.Message{
		Text: fmt.Sprintf("@%s liked your tweet", user.Username),
		UserID: tweet.UserID,
		AuthorID: user.ID,
		TweetID: tweet.ID,
	}

	if err := r.message.AddMessage(&message); err != nil {
		return err
	}

	return nil
}

func (r *TweetPostgres) UnlikeTweet(tweetId, userId uint) error {
	var tweet models.Tweet

	if err := r.db.Where("id = ?", tweetId).Preload("Likes").First(&tweet).Error; err != nil {
		return err
	}

	for _, user := range tweet.Likes {
		if user.ID == userId {
			if err := r.db.Model(&tweet).Association("Likes").Delete(user); err != nil {
				return err
			}
		}
	}

	return nil
}
