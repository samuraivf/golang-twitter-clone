package dto

type AddImageDto struct {
	Image string `json:"image" binding:"required"`
}
