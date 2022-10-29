package repository

import (
	"errors"
	"fmt"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPostgres struct {
	db *gorm.DB
	message Message
}

func NewUserPostgres(db *gorm.DB, message Message) *UserPostgres {
	return &UserPostgres{db, message}
}

func (r *UserPostgres) CreateUser(user dto.CreateUserDto) (uint, error) {
	userModel := models.User{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	result := r.db.Create(&userModel)

	if result.Error != nil {
		return 0, result.Error
	}

	return userModel.ID, nil
}

func (r *UserPostgres) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	result := r.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return user, nil
}

func (r *UserPostgres) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User
	result := r.db.Where("username = ?", username).Preload(clause.Associations).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return user, nil
}

func (r *UserPostgres) EditProfile(user dto.EditUserDto) error {
	userFromDb, err := r.GetUserByEmail(user.Email)

	if err != nil {
		return err
	}

	userFromDb.Name = user.Name
	userFromDb.Description = user.Description

	if err := r.db.Save(&userFromDb).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) AddImage(image []byte, userId uint) error {
	var user models.User
	err := r.db.First(&user, userId).Error

	if err != nil {
		return err
	}

	user.Image = image

	if err = r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) Subscribe(subscriberId, userId uint) error {
	var subscriber models.User
	var user models.User // the user being subscribed to

	if err := r.db.Where("id = ?", subscriberId).Preload("Subscriptions").First(&subscriber).Error; err != nil {
		return err
	}

	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return err
	}

	if err := r.db.Model(&subscriber).Association("Subscriptions").Append(&user); err != nil {
		return err
	}

	message := models.Message{
		Text: fmt.Sprintf("@%s subscribed to you", subscriber.Username),
		UserID: user.ID,
		AuthorID: subscriber.ID,
	}

	if err := r.message.AddMessage(&message); err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) Unsubscribe(subscriberId, userId uint) error {
	var subscriber models.User
	var user models.User // the user being subscribed to

	if err := r.db.Where("id = ?", subscriberId).Preload("Subscriptions").First(&subscriber).Error; err != nil {
		return err
	}

	if err := r.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := r.db.Model(&subscriber).Association("Subscriptions").Delete(&user); err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) GetUserMessages(userId uint) ([]*models.Message, error) {
	return r.message.GetUserMessages(userId)
}
