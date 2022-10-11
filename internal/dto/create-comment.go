package dto

type CreateCommentDto struct {
	Text string `json:"text"`
	UserID uint `json:"userId"`
	RepliedID uint `json:"repliedId"`
}