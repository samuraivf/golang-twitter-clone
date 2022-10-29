package models

type User struct {
	Model
	Username      string     `json:"username" gorm:"unique;not null;"`
	Name          string     `json:"name" gorm:"not null:"`
	Email         string     `json:"email" gorm:"not null;"`
	Password      string     `json:"-" gorm:"not null;"`
	Description   string     `json:"description" gorm:"not null"`
	Image         []byte     `json:"image"`
	Tweets        []Tweet    `json:"tweets"`
	LikedTweets   []*Tweet   `json:"likedTweets" gorm:"many2many:user_tweets;"`
	LikedComments []*Comment `json:"likedComments" gorm:"many2many:user_comments;"`
	MentionedIn   []*Tweet   `json:"mentionedIn" gorm:"many2many:user_mentionedIn;"`
	Messages 	  []Message  `json:"messages"`
	Subscriptions []*User    `json:"subscriptions" gorm:"many2many:user_subscriptions;"`
	Subscribers   []*User    `json:"subscribers" gorm:"many2many:user_subscriptions;foreignKey:ID;joinForeignKey:SubscriptionID;"`
}
