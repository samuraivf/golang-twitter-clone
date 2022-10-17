package dto

type UpdateCommentDto struct {
	Text      string `json:"text" binding:"required,min=1,max=10000"`
	CommentID uint `json:"commentId" binding:"required"`
}