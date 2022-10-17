package dto

type LoginDto struct {
	Username string `json:"username" binding:"required,min=3,max=15"`
	Password string `json:"password" binding:"required,min=8,max=30"`
}
