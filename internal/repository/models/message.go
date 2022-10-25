package models

type Message struct {
	Model
	Text   string `json:"text"`
	UserID uint   `json:"userId" gorm:"not null;"`
	AuthorID uint `json:"authorId" gorm:"not null;"`
	TweetID uint	`json:"tweetId" gorm:"default:null"`
}
