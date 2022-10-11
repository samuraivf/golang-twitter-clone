package handler

import (
	"net/http"

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
