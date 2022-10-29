package handler

import (
	"net/http"

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

	tag, err := h.services.Tag.GetTagByName(name)
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

	tag, err := h.services.Tag.GetTagByIdWithTweets(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tag)
}

func getTagName(c *gin.Context) string {
	name := c.Param("name")

	if name == "" {
		newErrorResponse(c, http.StatusBadRequest, errInvalidNameParam)
		return name
	}

	return name
}
