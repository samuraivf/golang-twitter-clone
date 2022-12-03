package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      string    `json:"text"`
	Likes     []UserId	`json:"likes" gorm:"foreignKey:CommentID"`
	TweetID   uint 		`json:"tweetId" gorm:"not null;"`
	UserID    uint 		`json:"userId" gorm:"not null;"`
}

type UserId struct {
	ID        uint `gorm:"primarykey" json:"id" binding:"required"`
	UserID    uint `json:"userId"`
	CommentID uint `json:"commentId"`
}
