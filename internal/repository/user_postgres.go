package repository

import (
	"errors"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db}
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
	result := r.db.Where("username = ?", username).First(&user)

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
