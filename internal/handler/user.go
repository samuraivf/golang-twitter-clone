package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
)

const (
	errEmptyParamUsername      = "empty param username"
	errCannotFindUser          = "cannot find a user"
	errCannotUpdateUserProfile = "cannot update user profile"
	errInvalidImage            = "invalid image"
)

func (h *Handler) getUserByUsername(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		newErrorResponse(c, http.StatusBadRequest, errEmptyParamUsername)
		return
	}

	user, err := h.services.User.GetUserByUsername(username)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errCannotFindUser)
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) editProfile(c *gin.Context) {
	var user dto.EditUserDto

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.User.EditProfile(user)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errCannotUpdateUserProfile)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) addImage(c *gin.Context) {
	data, err := getUserData(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var image dto.AddImageDto

	if err := c.BindJSON(&image); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidImage)
		return
	}

	if err := h.services.User.AddImage(image.Image, data.UserId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) Subscribe(c *gin.Context) {
	subscriberId := getUserId(c)
	userId := getIdParam(c)
	if subscriberId == 0 || userId == 0 {
		return
	}

	err := h.services.User.Subscribe(subscriberId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	subscriberId := getUserId(c)
	userId := getIdParam(c)
	if subscriberId == 0 || userId == 0 {
		return
	}

	err := h.services.User.Unsubscribe(subscriberId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *Handler) getUserMessages(c *gin.Context) {
	userId := getUserId(c)
	if userId == 0 {
		return
	}

	messages, err := h.services.User.GetUserMessages(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, messages)
}