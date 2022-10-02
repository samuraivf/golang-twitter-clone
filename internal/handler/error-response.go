package handler

import "github.com/gin-gonic/gin"

type ErrorMessage struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statueCode int, message string) {
	c.AbortWithStatusJSON(statueCode, ErrorMessage{message})
}
