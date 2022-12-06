package repository

import (
	"errors"

	"tag/internal/repo/models"
	
	"gorm.io/gorm"
)

var (
	ErrTagAlreadyExists = errors.New("error tag already exists")
)

type Tag interface {
	CreateTag(name string) (uint, error)
	GetTagByName(name string) (*models.Tag, error)
	AddTweet(tagId, tweetId uint) error
	DeleteTweet(tweetId uint) error
	GetTagById(id uint) (*models.Tag, error)
}

type TagPostgres struct {
	db *gorm.DB
}

func NewTagPostgres(db *gorm.DB) *TagPostgres {
	return &TagPostgres{db}
}

func (r *TagPostgres) CreateTag(name string) (uint, error) {
	tag := models.Tag{
		Name: name,
	}

	if err := r.db.Save(&tag).Error; err != nil {
		return 0, err
	}

	return tag.ID, nil
}

func (r *TagPostgres) GetTagByName(name string) (*models.Tag, error) {
	var tag *models.Tag

	if err := r.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *TagPostgres) AddTweet(tagId, tweetId uint) error {
	tweet := &models.TweetId{
		TweetID: tweetId,
		TagID: tagId,
	}

	if err := r.db.Save(&tweet).Error; err != nil {
		return err
	}

	return nil
}

func (r *TagPostgres) DeleteTweet(tweetId uint) error {
	return r.db.Where("tweet_id = ?", tweetId).Delete(&models.TweetId{}).Error
}

func (r *TagPostgres) GetTagById(id uint) (*models.Tag, error) {
	var tag *models.Tag

	if err := r.db.Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
