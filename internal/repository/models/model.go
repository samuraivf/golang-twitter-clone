package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
