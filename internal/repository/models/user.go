package models

type User struct {
	Model
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
	Tweets      []Tweet `json:"tweets"`
	LikedTweets []*Tweet `json:"likedTweets" gorm:"many2many:user_tweets;"`
	LikedComments []*Comment `json:"likedComments" gorm:"many2many:user_comments;"`
	MentionedIn []*Tweet `json:"mentionedIn" gorm:"many2many:user_mentionedIn;"`
}
