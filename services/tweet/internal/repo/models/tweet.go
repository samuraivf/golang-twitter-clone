package models

import "time"

type Tweet struct {
	ID        	   uint      `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Text           string    `json:"text"`
	Likes          []UserId   `json:"likes" gorm:"foreignKey:TweetID"`
	UserID         uint      `json:"userId"`
	Comments       []CommentId `json:"comments" gorm:"foreignKey:TweetID"`
	MentionedUsers []MentionedUserId   `json:"mentionedUsers" gorm:"foreignKey:TweetID"`
	Tags           []TagId    `json:"tags" gorm:"foreignKey:TweetID"`
}

type UserId struct {
	ID      uint `gorm:"primarykey" json:"id" binding:"required"`
	UserID  uint `json:"userId"`
	TweetID uint `json:"tweetId"`
}

type MentionedUserId struct {
	ID      uint `gorm:"primarykey" json:"id" binding:"required"`
	UserID  uint `json:"userId"`
	TweetID uint `json:"tweetId"`
}

type TagId struct {
	ID      uint `gorm:"primarykey" json:"id" binding:"required"`
	TagID   uint `json:"tagId"`
	TweetID uint `json:"tweetId"`
}

type CommentId struct {
	ID      uint `gorm:"primarykey" json:"id" binding:"required"`
	CommentID uint `json:"commentId"`
	TweetID uint `json:"tweetId"`
}