package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserId(c *gin.Context) uint {
	userData, err := getUserData(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 0
	}

	return userData.UserId
}

func getTweetId(c *gin.Context) uint {
	id := c.Param("id")

	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, errEmptyTweetIdParam)
		return 0
	}

	tweetIdUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 0
	}

	return uint(tweetIdUint)
}