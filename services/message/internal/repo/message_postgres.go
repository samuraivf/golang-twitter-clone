package repository

import (
	"message/internal/repo/models"
	
	"gorm.io/gorm"
)

type Message interface {
	AddMessage(message *models.Message) error
	GetUserMessages(userId uint) ([]*models.Message, error)
}

type MessagePostgres struct {
	db *gorm.DB
}

func NewMessagePostgres(db *gorm.DB) *MessagePostgres {
	return &MessagePostgres{db}
}

func (r *MessagePostgres) AddMessage(message *models.Message) error {
	if err := r.db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

func (r *MessagePostgres) GetUserMessages(userId uint) ([]*models.Message, error) {
	var userMessages []*models.Message

	if err := r.db.Where("user_id = ?", userId).Find(&userMessages).Error; err != nil {
		return nil, err
	}

	return userMessages, nil
}