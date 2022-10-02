package service

import (
	"errors"

	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo}
}

func (s *UserService) CreateUser(user dto.CreateUserDto) (uint, error) {
	passwordHash, err := generatePasswordHash(user.Password)

	if err != nil {
		return 0, err
	}

	user.Password = string(passwordHash)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *UserService) ValidateUser(username, password string) (*models.User, error) {
	user, err := s.GetUserByUsername(username)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *UserService) EditProfile(user dto.EditUserDto) error {
	return s.repo.EditProfile(user)
}

func (s *UserService) AddImage(image string, userId uint) error {
	return s.repo.AddImage([]byte(image), userId)
}
