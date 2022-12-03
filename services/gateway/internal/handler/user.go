package handler

import (
	"gateway/internal/services"
	"gateway/internal/dto"
	userService "user/proto"
	messageService "message/proto"

	"net/http"

	"github.com/gin-gonic/gin"
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

	userConnection := services.ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	user, err := userClient.GetUserByUsername(c, &userService.Username{
		Username: username,
	})

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

	userConnection := services.ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	_, err := userClient.EditProfile(c, &userService.EditUserDto{
		Name: user.Name,
		Email: user.Email,
		Description: user.Description,
	})

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

	userConnection := services.ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	if _, err := userClient.AddImage(c, &userService.AddImageParams{
		Image: image.Image,
		UserId: data.UserId,
	}); err != nil {
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

	userConnection := services.ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	_, err := userClient.Subscribe(c, &userService.SubscriberUser{
		SubscriberId: uint64(subscriberId),
		UserId: uint64(userId),
	})
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

	userConnection := services.ConnectUserGrpc()
	defer userConnection.Close()

	userClient := userService.NewUserClient(userConnection)

	_, err := userClient.Unsubscribe(c, &userService.SubscriberUser{
		SubscriberId: uint64(subscriberId),
		UserId: uint64(userId),
	})
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
	
	messageConnection := services.ConnectMessageGrpc()
	defer messageConnection.Close()

	messageClient := messageService.NewMessageClient(messageConnection)

	messages, err := messageClient.GetUserMessages(c, &messageService.UserId{
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, messages)
}
