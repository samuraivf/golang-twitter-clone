package dto

type UpdateTweetDto struct {
	Text    string   `json:"text" binding:"required,min=1,max=10000"`
	TweetID uint     `json:"tweetId" binding:"required"`
	Tags    []string `json:"tags"`
}
