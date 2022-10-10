package models

type User struct {
	Model
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
	Tweets []Tweet
	LikedTweets []*Tweet `json:"likedTweets" gorm:"many2many:user_tweets;"`
}
