package handler

import (
	"gateway/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"

	tagService "tag/proto"
	tweetService "tweet/proto"
)

const (
	errInvalidNameParam = "err invalid name param"
)

func (h *Handler) getTagByName(c *gin.Context) {
	name := getTagName(c)
	if name == "" {
		return
	}

	tagConnection := services.ConnectTagGrpc()
	defer tagConnection.Close()

	tagClient := tagService.NewTagClient(tagConnection)

	tag, err := tagClient.GetTagByName(c, &tagService.TagName{
		Name: name,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *Handler) getTagByIdWithTweets(c *gin.Context) {
	id := getIdParam(c)
	if id == 0 {
		return
	}

	tagConnection := services.ConnectTagGrpc()
	defer tagConnection.Close()

	tagClient := tagService.NewTagClient(tagConnection)

	tweetConnection := services.ConnectTweetGrpc()
	defer tweetConnection.Close()

	tweetClient := tweetService.NewTweetClient(tweetConnection)

	tag, err := tagClient.GetTagById(c, &tagService.TagId{
		TagId: uint64(id),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tweets, err := tweetClient.GetTweetsByTagId(c, &tweetService.TagId{
		TagId: uint64(id),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	type tagWithTweets struct {
		Id uint64 `json:"id"`
		Name string `json:"name"`
		CreatedAt *timestamppb.Timestamp `json:"createdAt"`
		UpdatedAt *timestamppb.Timestamp `json:"updatedAt"`
		Tweets *tweetService.Tweets `json:"tweets"`
	}

	c.JSON(http.StatusOK, tagWithTweets{
		Id: tag.Id,
		Name: tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
		Tweets: tweets,
	})
}

func getTagName(c *gin.Context) string {
	name := c.Param("name")

	if name == "" {
		newErrorResponse(c, http.StatusBadRequest, errInvalidNameParam)
		return name
	}

	return name
}
