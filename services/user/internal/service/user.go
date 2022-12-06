package service

import (
	"context"
	"errors"

	"user/dto"
	"user/internal/repo"
	"user/internal/repo/models"
	
	pb "user/proto"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService struct {
	pb.UnimplementedUserServer
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *pb.CreateUserDto) (*pb.UserId, error) {
	passwordHash, err := generatePasswordHash(user.Password)
	if err != nil {
		return new(pb.UserId), err
	}

	user.Password = string(passwordHash)

	userDto := dto.CreateUserDto{
		Username: user.Username,
		Password: user.Password,
		Name: user.Name,
		Email: user.Email,
	}

	userId, err := s.repo.CreateUser(userDto)
	if err != nil {
		return new(pb.UserId), err
	}

	return &pb.UserId{UserId: uint64(userId)}, nil
}

func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email *pb.Email) (*pb.UserData, error) {
	user, err := s.repo.GetUserByEmail(email.Email)
	if err != nil {
		return new(pb.UserData), err
	}

	return toUserData(user), nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username *pb.Username) (*pb.UserData, error) {
	user, err := s.repo.GetUserByUsername(username.Username)
	if err != nil {
		return new(pb.UserData), err
	}

	return toUserData(user), nil
}

func (s *UserService) GetUserById(ctx context.Context, userId *pb.UserId) (*pb.UserData, error) {
	user, err := s.repo.GetUserById(uint(userId.UserId))
	if err != nil {
		return new(pb.UserData), err
	}

	return toUserData(user), nil
}

func (s *UserService) getUserWithPassword(username string) (*models.User, error) {
	return s.repo.GetUserByUsername(username)
} 

func (s *UserService) ValidateUser(ctx context.Context, params *pb.ValidateUserParams) (*pb.UserData, error) {
	user, err := s.getUserWithPassword(params.Username)
	if err != nil {
		return new(pb.UserData), err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return new(pb.UserData), ErrInvalidPassword
	}

	return toUserData(user), nil
}

func (s *UserService) EditProfile(ctx context.Context, user *pb.EditUserDto) (*emptypb.Empty, error) {
	editUserDto := dto.EditUserDto{
		Name: user.Name,
		Email: user.Email,
		Description: user.Description,
	}

	return new(emptypb.Empty), s.repo.EditProfile(editUserDto)
}

func (s *UserService) AddImage(ctx context.Context, params *pb.AddImageParams) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.AddImage([]byte(params.Image), uint(params.UserId))
}

func (s *UserService) Subscribe(ctx context.Context, params *pb.SubscriberUser) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.Subscribe(uint(params.SubscriberId), uint(params.UserId))
}

func (s *UserService) Unsubscribe(ctx context.Context, params *pb.SubscriberUser) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.Unsubscribe(uint(params.SubscriberId), uint(params.UserId))
}

func toUserData(user *models.User) *pb.UserData {
	var subscribers []*pb.UserData
	var subscriptions []*pb.UserData

	for _, sub := range user.Subscribers {
		subscribers = append(subscribers, toUserData(sub))
	}

	for _, sub := range user.Subscriptions {
		subscriptions = append(subscriptions, toUserData(sub))
	}

	return &pb.UserData{
		Id: uint64(user.ID),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Username: user.Username,
		Name: user.Name,
		Email: user.Email,
		Description: user.Description,
		Image: user.Image,
		Subscribers: subscribers,
		Subscriptions: subscriptions,
	}
}