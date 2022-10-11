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
		RepliedID: commentDto.RepliedID,
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

	result := r.db.First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}