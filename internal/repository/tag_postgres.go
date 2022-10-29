package repository

import (
	"errors"

	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

var (
	ErrTagAlreadyExists = errors.New("error tag already exists")
)

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

	if err := r.db.Where("name = ?", name).Preload("Tweets").First(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *TagPostgres) AddTweet(tagName string, tweetId uint) error {
	tag, err := r.GetTagByName(tagName)
	if err != nil {
		return err
	}

	var tweet models.Tweet

	err = r.db.Where("id = ?", tweetId).First(&tweet).Error
	if err != nil {
		return err
	}

	
	if err := r.db.Model(&tag).Association("Tweets").Append(&tweet); err != nil {
		return err
	}

	return nil
}

func (r *TagPostgres) GetTagByIdWithTweets(id uint) (*models.Tag, error) {
	var tag *models.Tag

	if err := r.db.Where("id = ?", id).Preload("Tweets").First(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
