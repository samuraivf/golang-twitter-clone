package handler

import (
	"net/http"
	"strconv"

	"gateway/internal/dto"
	"gateway/internal/services"
	
	tweetService "tweet/proto"

	"github.com/gin-gonic/gin"
)

const (
	errEmptyTweetIdParam = "empty tweet id param"
	errEmptyUserIdParam  = "empty userId param"
)

// @Summary CreateTweet
// @Security ApiKeyAuth
// @Tags tweet
// @Description Create tweet
// @ID create-tweet
// @Accept json
// @Produce json
// @Param input body dto.CreateTweetDto true "create tweet data"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /tweet/create [post]
func (h *Handler) createTweet(c *gin.Context) {
	var tweetDto dto.CreateTweetDto

	if err := c.BindJSON(&tweetDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	tweetID, err := tweetClient.CreateTweet(c, &tweetService.CreateTweetDto{
		Text: tweetDto.Text,
		UserId: uint64(tweetDto.UserID),
		Tags: tweetDto.Tags,
		AuthorUsername: tweetDto.AuthorUsername,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"tweetId": uint(tweetID.TweetId),
	})
}

// @Summary UpdateTweet
// @Security ApiKeyAuth
// @Tags tweet
// @Description Update tweet
// @ID update-tweet
// @Accept json
// @Produce json
// @Param input body dto.UpdateTweetDto true "update tweet data"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /tweet/update [put]
func (h *Handler) updateTweet(c *gin.Context) {
	var tweetDto dto.UpdateTweetDto

	if err := c.BindJSON(&tweetDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	tweetID, err := tweetClient.UpdateTweet(c, &tweetService.UpdateTweetDto{
		Text: tweetDto.Text,
		TweetId: uint64(tweetDto.TweetID),
		Tags: tweetDto.Tags,
		AuthorUsername: tweetDto.AuthorUsername,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"tweetId": uint(tweetID.TweetId),
	})
}

// @Summary LikeTweet
// @Security ApiKeyAuth
// @Tags tweet
// @Description Like tweet
// @ID like-tweet
// @Produce json
// @Param id path string true "tweet id"
// @Success 200 {boolean} boolean true
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /tweet/like/{id} [get]
func (h *Handler) likeTweet(c *gin.Context) {
	userId := getUserId(c)
	tweetId := getIdParam(c)
	if userId == 0 || tweetId == 0 {
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	_, err := tweetClient.LikeTweet(c, &tweetService.TweetUserId{
		TweetId: uint64(tweetId),
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Summary UnlikeTweet
// @Security ApiKeyAuth
// @Tags tweet
// @Description Unlike tweet
// @ID unlike-tweet
// @Produce json
// @Param id path string true "tweet id"
// @Success 200 {boolean} boolean true
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /tweet/unlike/{id} [get]
func (h *Handler) unlikeTweet(c *gin.Context) {
	userId := getUserId(c)
	tweetId := getIdParam(c)
	if userId == 0 || tweetId == 0 {
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()
	_, err := tweetClient.UnlikeTweet(c, &tweetService.TweetUserId{
		TweetId: uint64(tweetId),
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Summary DeleteTweet
// @Security ApiKeyAuth
// @Tags tweet
// @Description Delete tweet by tweet id
// @ID delete-tweet
// @Produce json
// @Param id path string true "tweet id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /tweet/{id} [delete]
func (h *Handler) deleteTweet(c *gin.Context) {
	tweetId := getIdParam(c)
	if tweetId == 0 {
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	_, err := tweetClient.DeleteTweet(c, &tweetService.TweetId{
		TweetId: uint64(tweetId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweetId)
}

// @Summary GetTweetById
// @Security ApiKeyAuth
// @Tags tweet
// @Description Get tweet by ID
// @ID get-tweet-by-id
// @Produce json
// @Param id path string true "tweet id"
// @Success 200 {object} tweetService.TweetData
// @Failure 404 {object} ErrorMessage
// @Router /tweet/{id} [get]
func (h *Handler) getTweetById(c *gin.Context) {
	tweetId := getIdParam(c)
	if tweetId == 0 {
		return
	}

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	tweet, err := tweetClient.GetTweetById(c, &tweetService.TweetId{
		TweetId: uint64(tweetId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweet)
}

// @Summary GetUserTweets
// @Security ApiKeyAuth
// @Tags tweet
// @Description Get user tweets by userId
// @ID get-user-tweets
// @Produce json
// @Param userId path string true "user id"
// @Success 200 {object} tweetService.Tweets
// @Failure 400 {object} ErrorMessage
// @Failure 404 {object} ErrorMessage
// @Router /tweet/user-tweets/{userId} [get]
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

	tweetClient, closeTweet := services.GetTweetClient()
	defer closeTweet()

	tweets, err := tweetClient.GetUserTweets(c, &tweetService.UserIdParam{
		UserId: userIdUint,
	})
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweets)
}
