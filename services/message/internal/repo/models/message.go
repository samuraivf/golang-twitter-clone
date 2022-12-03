package models

import "time"

type Message struct {
	ID        	  uint      `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Text   string `json:"text"`
	UserID uint   `json:"userId" gorm:"not null;"`
	AuthorID uint `json:"authorId" gorm:"not null;"`
	TweetID uint  `json:"tweetId" gorm:"default:null"`
}
