package dto

type EditUserDto struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Description string `json:"description" binding:"required"`
}
