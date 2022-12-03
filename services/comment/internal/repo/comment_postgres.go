package repository

import (
	"context"
	"errors"
	"fmt"

	"comment/dto"
	"comment/internal/repo/models"
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

	result := r.db.Save(&comment)
	if result.Error != nil {
		return 0, result.Error
	}

	tweetConnection := ConnectTweetGrpc()
	defer tweetConnection.Close()

	tweetClient := tweetService.NewTweetClient(tweetConnection)

	ctx := context.Background()

	if _, err := tweetClient.AddComment(ctx, &tweetService.CommentTweetId{
		CommentId: uint64(comment.ID),
		TweetId: uint64(comment.TweetID),
	}); err != nil {
		return 0, err
	}

	messageConnection := ConnectMessageGrpc()
	defer messageConnection.Close()

	messageClient := messageService.NewMessageClient(messageConnection)

	message := messageService.MessageData{
		Text:     fmt.Sprintf("@%s commented your tweet", commentDto.Username),
		UserId:   uint64(commentDto.TweetAuthorID),
		AuthorId: uint64(commentDto.UserID),
		TweetId:  uint64(commentDto.TweetID),
	}

	if _, err := messageClient.AddMessage(ctx, &message); err != nil {
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
	tweetConnection := ConnectTweetGrpc()
	defer tweetConnection.Close()

	tweetClient := tweetService.NewTweetClient(tweetConnection)

	if _, err := tweetClient.DeleteComment(context.Background(), &tweetService.CommentId{
		CommentId: uint64(id),
	}); err != nil {
		return err
	}

	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentPostgres) LikeComment(commentId, userId uint, username string) error {
	var comment models.Comment

	if err := r.db.Where("id = ?", commentId).Preload("Likes").First(&comment).Error; err != nil {
		return err
	}

	for _, user := range comment.Likes {
		if user.ID == userId {
			return errors.New("user has already liked this tweet")
		}
	}

	if err := r.db.Model(&comment).Association("Likes").Append(&models.UserId{
		UserID: userId,
		CommentID: comment.ID,
	}); err != nil {
		return err
	}

	connection := ConnectMessageGrpc()
	defer connection.Close()

	messageClient := messageService.NewMessageClient(connection)

	message := messageService.MessageData{
		Text:     fmt.Sprintf("@%s liked your comment", username),
		UserId:   uint64(comment.UserID),
		AuthorId: uint64(userId),
		TweetId:  uint64(comment.TweetID),
	}

	if _, err := messageClient.AddMessage(context.Background(), &message); err != nil {
		return err
	}

	return nil
}

func (r *CommentPostgres) UnlikeComment(commentId, userId uint) error {
	var comment models.Comment

	if err := r.db.Where("id = ?", commentId).Preload("Likes").First(&comment).Error; err != nil {
		return err
	}

	for _, user := range comment.Likes {
		if user.UserID == userId {
			if err := r.db.Model(&comment).Association("Likes").Delete(user); err != nil {
				return err
			}

			if err := r.db.Delete(user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
