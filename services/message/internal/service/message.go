package service

import (
	"message/internal/repo/models"
	
	repository "message/internal/repo"
	pb "message/proto"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageService struct {
	pb.UnimplementedMessageServer
	repo repository.Message
}

func NewMessageService(repo repository.Message) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) AddMessage(ctx context.Context, message *pb.MessageData) (*emptypb.Empty, error) {
	messageModel := models.Message{
		ID: uint(message.Id),
		Text: message.Text,
		UserID: uint(message.UserId),
		AuthorID: uint(message.AuthorId),
		TweetID: uint(message.TweetId),
	}

	return new(emptypb.Empty), s.repo.AddMessage(&messageModel)
}

func (s *MessageService) GetUserMessages(ctx context.Context, userId *pb.UserId) (*pb.Messages, error) {
	messageModels, err := s.repo.GetUserMessages(uint(userId.UserId))
	if err != nil {
		return new(pb.Messages), err
	}

	var messages []*pb.MessageData

	for _, messageModel := range messageModels {
		message := &pb.MessageData{
			Id: uint64(messageModel.ID),
			CreatedAt: timestamppb.New(messageModel.CreatedAt),
			UpdatedAt: timestamppb.New(messageModel.UpdatedAt),
			Text: messageModel.Text,
			UserId: uint64(messageModel.UserID),
			AuthorId: uint64(messageModel.AuthorID),
			TweetId: uint64(messageModel.TweetID),
		}

		messages = append(messages, message)
	}

	return &pb.Messages{
		Messages: messages,
	}, nil
}