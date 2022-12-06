package handler

import (
	"net/http"
	
	"gateway/internal/services"
	
	tagService "tag/proto"
	tweetService "tweet/proto"
	
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/gin-gonic/gin"
)

const (
	errInvalidNameParam = "err invalid name param"
)

func (h *Handler) getTagByName(c *gin.Context) {
	name := getTagName(c)
	if name == "" {
		return
	}

	tagClient, closeTag := services.GetTagClient()
	defer closeTag()

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

	tagClient, closeTag := services.GetTagClient()
	defer closeTag()

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

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
