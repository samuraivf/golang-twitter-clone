package dto

type CreateTweetDto struct {
	Text   string   `json:"text" binding:"required,min=1,max=10000"`
	UserID uint     `json:"userId" binding:"required"`
	Tags   []string `json:"tags"`
	AuthorUsername string `json:"authorUsername"`
}
