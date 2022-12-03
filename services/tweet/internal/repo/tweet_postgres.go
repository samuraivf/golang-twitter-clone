package repository

import (
	"context"
	"errors"
	"fmt"

	messageService "message/proto"
	userService "user/proto"
	"tweet/dto"
	"tweet/internal/repo/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Tweet interface {
	CreateTweet(tweetDto dto.CreateTweetDto, mentionedUsers []string) (uint, error)
	GetTweetById(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]models.Tweet, error)
	GetTweetsByTagId(tagId uint) ([]models.Tweet, error)
	AddTag(tagId, tweetId uint) error
	AddComment(commentId, tweetId uint) error
	DeleteComment(commentId uint) error
	UpdateTweet(tweetDto dto.UpdateTweetDto, mentionedUsers []string) (uint, error)
	DeleteTweet(tweetId uint) error
	LikeTweet(tweetId, userId uint) error
	UnlikeTweet(tweetId, userId uint) error
}

var (
	errUserAlreadyLikedATweet = errors.New("user have already liked a tweet")
)

type TweetPostgres struct {
	db *gorm.DB
}

func NewTweetPostgres(db *gorm.DB) *TweetPostgres {
	return &TweetPostgres{db}
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

	if err := r.notifySubscribers(tweet.UserID, tweet.ID, tweetDto.AuthorUsername); err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

func (r *TweetPostgres) notifySubscribers(authorId, tweetId uint, authorUsername string) error {
	messageConnection := ConnectMessageGrpc()
	defer messageConnection.Close()

	messageClient := messageService.NewMessageClient(messageConnection)

	userConnection := ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	tweetAuthor, err := userClient.GetUserByUsername(
		context.Background(),
		&userService.Username{
			Username: authorUsername,
		},
	)

	if err != nil {
		return err
	}

	for _, subscriber := range tweetAuthor.Subscribers {
		message := messageService.MessageData{
			Text: fmt.Sprintf("@%s posted a new tweet", authorUsername),
			UserId: subscriber.Id,
			AuthorId: tweetAuthor.Id,
			TweetId: uint64(tweetId),
		}

		if _, err := messageClient.AddMessage(context.Background(), &message); err != nil {
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

func (r *TweetPostgres) AddTag(tagId, tweetId uint) error {
	tag := &models.TagId{
		TagID: tagId,
		TweetID: tweetId,
	}

	if err := r.db.Save(&tag).Error; err != nil {
		return err
	}

	return nil
}

func (r *TweetPostgres) AddComment(commentId, tweetId uint) error {
	comment := &models.CommentId{
		CommentID: commentId,
		TweetID: tweetId,
	}

	if err := r.db.Save(&comment).Error; err != nil {
		return err
	}

	fmt.Println(comment)

	return nil
}

func (r *TweetPostgres) DeleteComment(commentId uint) error {
	return r.db.Delete(&models.CommentId{}, commentId).Error
}

func (r *TweetPostgres) addMentions(tweet *models.Tweet, mentionedUsers []string) error {
	var tweetAuthor *userService.UserData

	messageConnection := ConnectMessageGrpc()
	defer messageConnection.Close()

	messageClient := messageService.NewMessageClient(messageConnection)

	userConnection := ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	tweetAuthor, err := userClient.GetUserById(context.Background(), &userService.UserId{
		UserId: uint64(tweet.UserID),
	})

	if err != nil {
		return err
	}

	for _, user := range mentionedUsers {
		userFromDB, err := userClient.GetUserByUsername(context.Background(), &userService.Username{
			Username: user,
		})

		if err != nil {
			return err
		}

		if err := r.db.Model(&tweet).Association("MentionedUsers").Append(&models.MentionedUserId{
			UserID: uint(userFromDB.Id),
			TweetID: tweet.ID,
		}); err != nil {
			return err
		}

		message := messageService.MessageData{
			Text: fmt.Sprintf("You was metioned by @%s", tweetAuthor.Username),
			UserId: userFromDB.Id,
			AuthorId: uint64(tweet.UserID),
			TweetId: uint64(tweet.ID),
		}

		if _, err := messageClient.AddMessage(context.Background(), &message); err != nil {
			return err
		}
	}

	return nil
}

func (r *TweetPostgres) DeleteTweet(tweetId uint) error {
	err := r.db.Exec("DELETE FROM tag_ids WHERE tweet_id = ?", tweetId).Error

	if err != nil {
		return err
	}

	err = r.db.Exec("DELETE FROM comment_ids WHERE tweet_id = ?", tweetId).Error

	if err != nil {
		return err
	}

	err = r.db.Exec("DELETE FROM user_ids WHERE tweet_id = ?", tweetId).Error

	if err != nil {
		return err
	}

	return r.db.Delete(&models.Tweet{}, tweetId).Error
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

func (r *TweetPostgres) GetTweetsByTagId(tagId uint) ([]models.Tweet, error) {
	var tweets []models.Tweet

	result := r.db.Exec(
		"SELECT * FROM tweets WHERE tweets.id IN (SELECT tweet_id FROM tag_ids WHERE tag_id = ?);",
		tagId,
	).Find(&tweets)

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
		if user.UserID == userId {
			return errUserAlreadyLikedATweet
		}
	}

	if err := r.db.Model(&tweet).Association("Likes").Append(&models.UserId{
		UserID: userId,
		TweetID: tweetId,
	}); err != nil {
		return err
	}

	messageConnection := ConnectMessageGrpc()
	defer messageConnection.Close()

	messageClient := messageService.NewMessageClient(messageConnection)

	userConnection := ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	user, err := userClient.GetUserById(context.Background(), &userService.UserId{
		UserId: uint64(userId),
	})

	if err != nil {
		return err
	}

	message := messageService.MessageData{
		Text: fmt.Sprintf("@%s liked your tweet", user.Username),
		UserId: uint64(tweet.UserID),
		AuthorId: user.Id,
		TweetId: uint64(tweet.ID),
	}

	if _, err := messageClient.AddMessage(context.Background(), &message); err != nil {
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
		if user.UserID == userId {
			if err := r.db.Model(&tweet).Association("Likes").Delete(user); err != nil {
				return err
			}

			if err := r.db.Delete(user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
