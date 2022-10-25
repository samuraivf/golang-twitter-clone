package models

type Tag struct {
	Model
	Name   string   `json:"string" gorm:"unique;not null;"`
	Tweets []*Tweet `json:"tweets" gorm:"many2many:tag_tweets;"`
}
