package dto

type CreateUserDto struct {
	LoginDto
	Name  string `json:"name" binding:"required,min=3,max=20"`
	Email string `json:"email" binding:"required,email"`
}
