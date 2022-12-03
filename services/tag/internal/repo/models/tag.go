package models

import "time"

type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name   string   `json:"string" gorm:"unique;not null;"`
	Tweets []TweetId `json:"tweets" gorm:"foreignKey:TagID"`
}

type TweetId struct {
	ID    uint `gorm:"primarykey" json:"id" binding:"required"`
	TweetID uint `json:"tweetId"`
	TagID uint `json:"tagId"`
}
