package models

type Comment struct {
	Model
	Text    string  `json:"text"`
	Likes   []*User `json:"likes" gorm:"many2many:user_comments;"`
	TweetID uint    `json:"tweetId" gorm:"unique;not null;"`
	UserID  uint    `json:"userId" gorm:"unique;not null;"`
}
