package service

import (
	"context"
	
	"comment/dto"
	"comment/internal/repo"
	
	pb "comment/proto"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommentService struct {
	pb.UnimplementedCommentServer
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(ctx context.Context, commentDto *pb.CreateCommentDto) (*pb.CommentId, error) {
	comment := dto.CreateCommentDto{
		Text: commentDto.Text,
		UserID: uint(commentDto.UserId),
		TweetID: uint(commentDto.TweetId),
		Username: commentDto.Username,
		TweetAuthorID: uint(commentDto.TweetAuthorId),
	}

	id, err := s.repo.CreateComment(comment)

	if err != nil {
		return new(pb.CommentId), err
	}

	return &pb.CommentId{CommentId: uint64(id)}, nil
}

func (s *CommentService) GetCommentById(ctx context.Context, id *pb.CommentId) (*pb.CommentData, error) {
	commentModel, err := s.repo.GetCommentById(uint(id.CommentId))

	if err != nil {
		return new(pb.CommentData), err
	}

	var likes []*pb.UserId

	for _, userId := range commentModel.Likes {
		likes = append(likes, &pb.UserId{
			UserId: uint64(userId.ID),
			CommentId: uint64(userId.CommentID),
		})
	}

	return &pb.CommentData{
		Id: uint64(commentModel.ID),
		CreatedAt: timestamppb.New(commentModel.CreatedAt),
		UpdatedAt: timestamppb.New(commentModel.UpdatedAt),
		Text: commentModel.Text,
		Likes: likes,
		TweetId: uint64(commentModel.TweetID),
		UserId: uint64(commentModel.UserID),
	}, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, commentDto *pb.UpdateCommentDto) (*pb.CommentId, error) {
	comment := dto.UpdateCommentDto{
		Text: commentDto.Text,
		CommentID: uint(commentDto.CommentId),
	}

	id, err := s.repo.UpdateComment(comment)

	if err != nil {
		return new(pb.CommentId), err
	}

	return &pb.CommentId{CommentId: uint64(id)}, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id *pb.CommentId) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.DeleteComment(uint(id.CommentId))
}

func (s *CommentService) LikeComment(ctx context.Context, commentUser *pb.CommentUser) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.LikeComment(uint(commentUser.CommentId), uint(commentUser.UserId), commentUser.Username)
}

func (s *CommentService) UnlikeComment(ctx context.Context, commentUser *pb.CommentUser) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.UnlikeComment(uint(commentUser.CommentId), uint(commentUser.UserId))
}
