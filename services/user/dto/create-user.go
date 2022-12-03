package dto

type CreateUserDto struct {
	Username string `json:"username" binding:"required,min=3,max=15"`
	Password string `json:"password" binding:"required,min=8,max=30"`
	Name  string `json:"name" binding:"required,min=3,max=20"`
	Email string `json:"email" binding:"required,email"`
}