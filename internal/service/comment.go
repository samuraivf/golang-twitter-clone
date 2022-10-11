package service

import (
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo}
}

func (s *CommentService) CreateComment(commentDto dto.CreateCommentDto) (uint, error) {
	return s.repo.CreateComment(commentDto)
}

func (s *CommentService) GetCommentById(id uint) (*models.Comment, error) {
	return s.repo.GetCommentById(id)
}