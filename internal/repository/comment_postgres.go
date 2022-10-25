package repository

import (
	"fmt"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

type CommentPostgres struct {
	db      *gorm.DB
	message Message
}

func NewCommentPostgres(db *gorm.DB, message Message) *CommentPostgres {
	return &CommentPostgres{db, message}
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

	var user *models.User
	var tweet *models.Tweet

	if err := r.db.Select("username").First(&user, commentDto.UserID).Error; err != nil {
		return 0, err
	}
	if err := r.db.Select("user_id").First(&tweet, commentDto.TweetID).Error; err != nil {
		return 0, err
	}

	message := models.Message{
		Text:     fmt.Sprintf("@%s commented your tweet", user.Username),
		UserID:   tweet.UserID,
		AuthorID: commentDto.UserID,
		TweetID:  commentDto.TweetID,
	}

	if err := r.message.AddMessage(&message); err != nil {
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
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentPostgres) LikeComment(commentId, userId uint) error {
	var comment models.Comment

	if err := r.db.Where("id = ?", commentId).Preload("Likes").First(&comment).Error; err != nil {
		return err
	}

	for _, user := range comment.Likes {
		if user.ID == userId {
			return errUserAlreadyLikedATweet
		}
	}

	var user models.User

	if err := r.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := r.db.Model(&comment).Association("Likes").Append(&user); err != nil {
		return err
	}

	message := models.Message{
		Text:     fmt.Sprintf("@%s liked your comment", user.Username),
		UserID:   comment.UserID,
		AuthorID: user.ID,
		TweetID:  comment.TweetID,
	}

	if err := r.message.AddMessage(&message); err != nil {
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
		if user.ID == userId {
			if err := r.db.Model(&comment).Association("Likes").Delete(user); err != nil {
				return err
			}
		}
	}

	return nil
}
