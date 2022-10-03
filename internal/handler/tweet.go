package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
)

const (
	errEmptyTweetIdParam = "empty tweet id param"
	errEmptyUserIdParam = "empty userId param"
)

func (h *Handler) createTweet(c *gin.Context) {
	var tweetDto dto.CreateTweetDto

	if err := c.BindJSON(&tweetDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	tweetID, err := h.services.Tweet.CreateTweet(tweetDto)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"tweetId": tweetID,
	})
}

func (h *Handler) getTweetById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, errEmptyTweetIdParam)
		return
	}

	tweet, err := h.services.Tweet.GetTweetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (h *Handler) getUserTweets(c *gin.Context) {
	userId := c.Param("userId")

	if userId == "" {
		newErrorResponse(c, http.StatusBadRequest, errEmptyUserIdParam)
		return
	} 

	userIdUint, err := strconv.ParseUint(userId, 10, 64)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tweets, err := h.services.Tweet.GetUserTweets(uint(userIdUint))

	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweets)
}