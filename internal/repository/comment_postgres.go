package repository

import (
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

type CommentPostgres struct {
	db *gorm.DB
}

func NewCommentPostgres(db *gorm.DB) *CommentPostgres {
	return &CommentPostgres{db}
}

func (r *CommentPostgres) CreateComment(commentDto dto.CreateCommentDto) (uint, error) {
	comment := models.Comment{
		Text: commentDto.Text,
		TweetID: commentDto.TweetID,
		UserID: commentDto.UserID,
	}

	result := r.db.Save(&comment)
	if result.Error != nil {
		return 0, result.Error
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