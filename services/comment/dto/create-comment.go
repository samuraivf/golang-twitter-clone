package dto

type CreateCommentDto struct {
	Text    string `json:"text" binding:"required,min=1,max=10000"`
	UserID  uint   `json:"userId" binding:"required"`
	Username string `json:"username" binding:"required"`
	TweetID uint   `json:"tweetId" binding:"required"`
	TweetAuthorID uint `json:"tweetAuthorId" binding:"required"`
}
