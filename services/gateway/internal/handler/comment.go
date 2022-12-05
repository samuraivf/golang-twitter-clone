package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gateway/internal/dto"
	"gateway/internal/services"
	commentService "comment/proto"
)

const (
	errEmptyCommentIdParam = "error empty comment id param"
)

func (h *Handler) createComment(c *gin.Context) {
	var commentDto dto.CreateCommentDto

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	if err := c.BindJSON(&commentDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	tweetId, err := commentClient.CreateComment(c, &commentService.CreateCommentDto{
		Text: commentDto.Text,
		UserId: uint64(commentDto.UserID),
		TweetId: uint64(commentDto.TweetID),
		Username: commentDto.Username,
		TweetAuthorId: uint64(commentDto.TweetAuthorID),
	})
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

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	comment, err := commentClient.GetCommentById(c, &commentService.CommentId{
		CommentId: uint64(id),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *Handler) updateComment(c *gin.Context) {
	var commentDto dto.UpdateCommentDto

	if err := c.BindJSON(&commentDto); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	id, err := commentClient.UpdateComment(c, &commentService.UpdateCommentDto{
		Text: commentDto.Text,
		CommentId: uint64(commentDto.CommentID),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (h *Handler) deleteComment(c *gin.Context) {
	id := getCommentId(c)
	if id == 0 {
		return
	}

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	_, err := commentClient.DeleteComment(c, &commentService.CommentId{
		CommentId: uint64(id),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (h *Handler) likeComment(c *gin.Context) {
	userId := getUserId(c)
	commentId := getCommentId(c)
	if userId == 0 || commentId == 0 {
		return
	}

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	_, err := commentClient.LikeComment(c, &commentService.CommentUser{
		CommentId: uint64(commentId),
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *Handler) unlikeComment(c *gin.Context) {
	userId := getUserId(c)
	commentId := getCommentId(c)
	if userId == 0 || commentId == 0 {
		return
	}

	commentClient, closeComment := services.GetCommentClient()
	defer closeComment()

	_, err := commentClient.UnlikeComment(c, &commentService.CommentUser{
		CommentId: uint64(commentId),
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
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
