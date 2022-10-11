package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
)

const (
	errEmptyCommentIdParam = "error empty comment id param"
)

func (h *Handler) createComment(c *gin.Context) {
	var commentDto dto.CreateCommentDto

	if err := c.BindJSON(&commentDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	tweetId, err := h.services.Comment.CreateComment(commentDto)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tweetId)
}

func (h *Handler) getCommentById(c *gin.Context) {
	id := getCommentId(c)
	if id == 0 {
		return
	}

	comment, err := h.services.Comment.GetCommentById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, comment)
}

func getCommentId(c *gin.Context) uint {
	id := c.Param("id")

	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, errEmptyCommentIdParam)
		return 0
	}

	commentIdUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 0
	}

	return uint(commentIdUint)
}