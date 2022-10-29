package service

import (
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo}
}

func (s *TagService) GetTagByName(name string) (*models.Tag, error) {
	return s.repo.GetTagByName(name)
}

func (s *TagService) GetTagByIdWithTweets(id uint) (*models.Tag, error) {
	return s.repo.GetTagByIdWithTweets(id)
}
