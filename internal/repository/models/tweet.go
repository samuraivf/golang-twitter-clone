package models

type Tweet struct {
	Model
	Text string `json:"text"`
	Likes uint `json:"likes"`
	Retweets uint `json:"retweets"`
	UserID uint `json:"userId"`
}