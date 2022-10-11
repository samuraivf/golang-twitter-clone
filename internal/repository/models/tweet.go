package models

type Tweet struct {
	Model
	Text     string  `json:"text"`
	Likes    []*User `json:"likes" gorm:"many2many:user_tweets;"`
	Retweets uint    `json:"retweets"`
	UserID   uint    `json:"userId"`
	Comments []Comment `json:"comments" gorm:"foreignKey:RepliedID"`
}
