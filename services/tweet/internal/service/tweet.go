package service

import (
	"context"
	"regexp"
	"strings"

	"tweet/dto"
	"tweet/internal/connections"
	"tweet/internal/repo"
	"tweet/internal/repo/models"

	pb "tweet/proto"
	tagService "tag/proto"
	
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TweetService struct {
	pb.UnimplementedTweetServer
	repo repository.Tweet
}

func NewTweetService(repo repository.Tweet) *TweetService {
	return &TweetService{repo: repo}
}

func (s *TweetService) CreateTweet(ctx context.Context, tweetDto *pb.CreateTweetDto) (*pb.TweetId, error) {
	var mentionedUsers []string

	for _, letter := range strings.Split(tweetDto.Text, " ") {
		if letter[0] == byte('@') {
			mentionedUsers = append(mentionedUsers, letter[1:])
		}
	}

	tweetData := dto.CreateTweetDto{
		Text: tweetDto.Text,
		UserID: uint(tweetDto.UserId),
		Tags: tweetDto.Tags,
		AuthorUsername: tweetDto.AuthorUsername,
	}

	tweetId, err := s.repo.CreateTweet(tweetData, mentionedUsers)
	if err != nil {
		return new(pb.TweetId), err
	}

	return s.addTags(tweetDto.Tags, tweetId)
}

func (s *TweetService) GetTweetById(ctx context.Context, id *pb.TweetId) (*pb.TweetData, error) {
	tweet, err := s.repo.GetTweetById(uint(id.TweetId))
	if err != nil {
		return new(pb.TweetData), err
	}

	return toTweetData(tweet), nil
}

func (s *TweetService) GetUserTweets(ctx context.Context, userId *pb.UserIdParam) (*pb.Tweets, error) {
	tweets, err := s.repo.GetUserTweets(uint(userId.UserId))
	if err != nil {
		return new(pb.Tweets), err
	}

	var pbTweets []*pb.TweetData

	for _, tweet := range tweets {
		pbTweets = append(pbTweets, toTweetData(&tweet))
	}

	return &pb.Tweets{
		Tweets: pbTweets,
	}, nil
}

func (s *TweetService) UpdateTweet(ctx context.Context, tweetDto *pb.UpdateTweetDto) (*pb.TweetId, error) {
	var mentionedUsers []string

	for _, letter := range strings.Split(tweetDto.Text, " ") {
		if letter[0] == byte('@') {
			mentionedUsers = append(mentionedUsers, letter[1:])
		}
	}

	tweetData := dto.UpdateTweetDto{
		Text: tweetDto.Text,
		TweetID: uint(tweetDto.TweetId),
		Tags: tweetDto.Tags,
		AuthorUsername: tweetDto.AuthorUsername,
	}

	tweetId, err := s.repo.UpdateTweet(tweetData, mentionedUsers)
	if err != nil {
		return new(pb.TweetId), err
	}

	return &pb.TweetId{TweetId: uint64(tweetId)}, nil
}

func (s *TweetService) DeleteTweet(ctx context.Context, tweetId *pb.TweetId) (*emptypb.Empty, error) {
	err := s.repo.DeleteTweet(uint(tweetId.TweetId))
	if err != nil {
		return new(emptypb.Empty), err
	}

	tagClient, closeTag := connections.GetTagClient()
	defer closeTag()

	_, err = tagClient.DeleteTweet(ctx, &tagService.TweetIdParam{
		TweetId: tweetId.TweetId,
	})

	return new(emptypb.Empty), err
}

func (s *TweetService) LikeTweet(ctx context.Context, params *pb.TweetUserId) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.LikeTweet(uint(params.TweetId), uint(params.UserId))
}

func (s *TweetService) UnlikeTweet(ctx context.Context, params *pb.TweetUserId) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.UnlikeTweet(uint(params.TweetId), uint(params.UserId))
}

func (s *TweetService) AddComment(ctx context.Context, params *pb.CommentTweetId) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.AddComment(uint(params.CommentId), uint(params.TweetId))
}

func (s *TweetService) DeleteComment(ctx context.Context, params *pb.CommentId) (*emptypb.Empty, error) {
	return new(emptypb.Empty), s.repo.DeleteComment(uint(params.CommentId))
}

func (s *TweetService) GetTweetsByTagId(ctx context.Context, tagId *pb.TagId) (*pb.Tweets, error) {
	tweets, err := s.repo.GetTweetsByTagId(uint(tagId.TagId))
	if err != nil {
		return new(pb.Tweets), err
	}

	var pbTweets []*pb.TweetData
	
	for _, tweet := range tweets {
		pbTweets = append(pbTweets, toTweetData(&tweet))
	}

	return &pb.Tweets{
		Tweets: pbTweets,
	}, nil
}

func (s *TweetService) addTags(tags []string, tweetId uint) (*pb.TweetId, error) {
	tagClient, closeTag := connections.GetTagClient()
	defer closeTag()

	for _, tag := range tags {
		r, err := regexp.Compile("^[a-z0-9]+$")
		if err != nil {
			return new(pb.TweetId), err
		}

		if r.Match([]byte(tag)) {
			ctx := context.Background()

			tagData, err := tagClient.GetTagByName(ctx, &tagService.TagName{
				Name: tag,
			})

			var tagId uint64

			if err != nil {
				pbTagId, err := tagClient.CreateTag(ctx, &tagService.TagName{Name: tag})
				if err != nil {
					return new(pb.TweetId), err
				}

				tagId = pbTagId.TagId
			} else {
				tagId = tagData.Id
			}

			_, err = tagClient.AddTweet(ctx, &tagService.AddTweetParams{
				TagId: tagId,
				TweetId: uint64(tweetId),
			})
			if err != nil {
				return new(pb.TweetId), err
			}

			err = s.repo.AddTag(uint(tagId), tweetId)
			if err != nil {
				return new(pb.TweetId), err
			}
		}
	}

	return &pb.TweetId{TweetId: uint64(tweetId)}, nil
}

func toTweetData(tweet *models.Tweet) (*pb.TweetData) {
	var (
		likes []*pb.UserId
		comments []*pb.CommentId
		mentionedUsers []*pb.UserId
		tags []*pb.TagId
	)

	for _, userId := range tweet.Likes {
		likes = append(likes, &pb.UserId{
			UserId: uint64(userId.ID),
			TweetId: uint64(userId.TweetID),
		})
	}

	for _, commentId := range tweet.Comments{
		comments = append(comments, &pb.CommentId{
			CommentId: uint64(commentId.ID),
			TweetId: uint64(commentId.TweetID),
		})
	}

	for _, userId := range tweet.MentionedUsers {
		mentionedUsers = append(mentionedUsers, &pb.UserId{
			UserId: uint64(userId.ID),
			TweetId: uint64(userId.TweetID),
		})
	}

	for _, tagId := range tweet.Tags {
		tags = append(tags, &pb.TagId{
			TagId: uint64(tagId.ID),
			TweetId: uint64(tagId.TweetID),
		})
	}


	return &pb.TweetData{
		Id: uint64(tweet.ID),
		CreatedAt: timestamppb.New(tweet.CreatedAt),
		UpdatedAt: timestamppb.New(tweet.UpdatedAt),
		Text: tweet.Text,
		Likes: likes,
		UserId: uint64(tweet.UserID),
		Comments: comments,
		MentionedUsers: mentionedUsers,
		Tags: tags,
	}
}