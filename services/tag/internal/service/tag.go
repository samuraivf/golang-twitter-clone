package service

import (
	"context"

	"tag/internal/repo"
	"tag/internal/repo/models"
	pb "tag/proto"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TagService struct {
	pb.UnimplementedTagServer
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(ctx context.Context, name *pb.TagName) (*pb.TagId, error) {
	id, err := s.repo.CreateTag(name.Name)

	if err != nil {
		return new(pb.TagId), err
	}

	return &pb.TagId{TagId: uint64(id)}, nil
}

func (s *TagService) GetTagByName(ctx context.Context, name *pb.TagName) (*pb.TagData, error) {
	tagModel, err := s.repo.GetTagByName(name.Name)

	if err != nil {
		return new(pb.TagData), err
	}

	return toTagData(tagModel), nil
}

func (s *TagService) GetTagById(ctx context.Context, id *pb.TagId) (*pb.TagData, error) {
	tagModel, err := s.repo.GetTagById(uint(id.TagId))

	if err != nil {
		return new(pb.TagData), err
	}

	return toTagData(tagModel), nil
}

func (s *TagService) AddTweet(ctx context.Context, params *pb.AddTweetParams) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.AddTweet(uint(params.TagId), uint(params.TweetId))
}

func (s *TagService) DeleteTweet(ctx context.Context, tweetId *pb.TweetIdParam) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.DeleteTweet(uint(tweetId.TweetId))
}

func toTagData(tag *models.Tag) (*pb.TagData) {
	var tweets []*pb.TweetId

	for _, tweetId := range tag.Tweets {
		tweets = append(tweets, &pb.TweetId{
			TweetId: uint64(tweetId.ID),
			TagId: uint64(tag.ID),
		})
	}

	return &pb.TagData{
		Id: uint64(tag.ID),
		CreatedAt: timestamppb.New(tag.CreatedAt),
		UpdatedAt: timestamppb.New(tag.UpdatedAt),
		Name: tag.Name,
		Tweets: tweets,
	}
}
