package repository

import (
	"context"
	"errors"
	"fmt"

	"comment/dto"
	"comment/internal/repo/models"
	"comment/internal/connections"
	
	messageService "message/proto"
	tweetService "tweet/proto"

	"gorm.io/gorm"
)

type Comment interface {
	CreateComment(commentDto dto.CreateCommentDto) (uint, error)
	GetCommentById(id uint) (*models.Comment, error)
	UpdateComment(commentDto dto.UpdateCommentDto) (uint, error)
	DeleteComment(id uint) error
	LikeComment(commentId, userId uint, username string) error
	UnlikeComment(commentId, userId uint) error
}

type CommentPostgres struct {
	db *gorm.DB
}

func NewCommentPostgres(db *gorm.DB) *CommentPostgres {
	return &CommentPostgres{db}
}

func (r *CommentPostgres) CreateComment(commentDto dto.CreateCommentDto) (uint, error) {
	comment := models.Comment{
		Text:    commentDto.Text,
		TweetID: commentDto.TweetID,
		UserID:  commentDto.UserID,
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	err := tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tweetClient, closeTweet := connections.GetTweetClient()
	defer closeTweet()

	ctx := context.Background()

	_, err = tweetClient.AddComment(ctx, &tweetService.CommentTweetId{
		CommentId: uint64(comment.ID),
		TweetId: uint64(comment.TweetID),
	})
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	messageClient, closeMessage := connections.GetMessageClient()
	defer closeMessage()

	message := messageService.MessageData{
		Text:     fmt.Sprintf("@%s commented your tweet", commentDto.Username),
		UserId:   uint64(commentDto.TweetAuthorID),
		AuthorId: uint64(commentDto.UserID),
		TweetId:  uint64(commentDto.TweetID),
	}

	_, err = messageClient.AddMessage(ctx, &message)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit().Error
	if err != nil {
		return 0, err
	}
	
	return comment.ID, nil
}

func (r *CommentPostgres) GetCommentById(id uint) (*models.Comment, error) {
	var comment *models.Comment

	result := r.db.Where("id = ?", id).Preload("Likes").First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (r *CommentPostgres) UpdateComment(commentDto dto.UpdateCommentDto) (uint, error) {
	var comment *models.Comment

	result := r.db.First(&comment, commentDto.CommentID)
	if result.Error != nil {
		return 0, result.Error
	}

	comment.Text = commentDto.Text

	result = r.db.Save(&comment)
	if result.Error != nil {
		return 0, result.Error
	}

	return comment.ID, nil
}

func (r *CommentPostgres) DeleteComment(id uint) error {
	tweetClient, closeTweet := connections.GetTweetClient()
	defer closeTweet()

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	_, err := tweetClient.DeleteComment(context.Background(), &tweetService.CommentId{
		CommentId: uint64(id),
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&models.Comment{}, id).Error
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *CommentPostgres) LikeComment(commentId, userId uint, username string) error {
	var comment models.Comment

	err := r.db.Where("id = ?", commentId).Preload("Likes").First(&comment).Error
	if err != nil {
		return err
	}

	for _, user := range comment.Likes {
		if user.ID == userId {
			return errors.New("user has already liked this tweet")
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err = tx.Model(&comment).Association("Likes").Append(&models.UserId{
		UserID: userId,
		CommentID: comment.ID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	messageClient, closeMessage := connections.GetMessageClient()
	defer closeMessage()

	message := messageService.MessageData{
		Text:     fmt.Sprintf("@%s liked your comment", username),
		UserId:   uint64(comment.UserID),
		AuthorId: uint64(userId),
		TweetId:  uint64(comment.TweetID),
	}

	_, err = messageClient.AddMessage(context.Background(), &message)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CommentPostgres) UnlikeComment(commentId, userId uint) error {
	var comment models.Comment

	err := r.db.Where("id = ?", commentId).Preload("Likes").First(&comment).Error
	if err != nil {
		return err
	}

	for _, user := range comment.Likes {
		if user.UserID == userId {
			tx := r.db.Begin()
			if tx.Error != nil {
				return tx.Error
			}

			err := tx.Model(&comment).Association("Likes").Delete(user)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = tx.Delete(user).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			return tx.Commit().Error
		}
	}

	return nil
}
