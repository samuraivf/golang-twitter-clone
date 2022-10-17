package dto

type EditUserDto struct {
	Name        string `json:"name" binding:"required,min=3,max=20"`
	Email       string `json:"email" binding:"required,email"`
	Description string `json:"description" binding:"required,min=0,max=1000"`
}
