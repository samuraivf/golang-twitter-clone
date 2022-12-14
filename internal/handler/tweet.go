package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
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

	tweetID, err := h.services.Tweet.CreateTweet(tweetDto)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"tweetId": tweetID,
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

	tweetID, err := h.services.Tweet.UpdateTweet(tweetDto)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"tweetId": tweetID,
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

	err := h.services.Tweet.LikeTweet(tweetId, userId)

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

	err := h.services.Tweet.UnlikeTweet(tweetId, userId)

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

	err := h.services.Tweet.DeleteTweet(tweetId)

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
// @Success 200 {object} models.Tweet
// @Failure 404 {object} ErrorMessage
// @Router /tweet/{id} [get]
func (h *Handler) getTweetById(c *gin.Context) {
	tweetId := getIdParam(c)

	if tweetId == 0 {
		return
	}

	tweet, err := h.services.Tweet.GetTweetById(tweetId)

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
// @Success 200 {array} models.Tweet
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

	tweets, err := h.services.Tweet.GetUserTweets(uint(userIdUint))

	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweets)
}
