package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/service"
)

func (h *Handler) signUp(c *gin.Context) {
	var user dto.CreateUserDto

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if _, err := h.services.User.GetUserByEmail(user.Email); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "user have already exists")
		return
	}

	if _, err := h.services.User.GetUserByUsername(user.Username); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "user have already exists")
		return
	}

	id, err := h.services.User.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var user dto.LoginDto

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	userData, err := h.services.User.ValidateUser(user.Username, user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.createTokens(c, userData.Username, userData.ID)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cannot logout: invalid refreshToken cookie")
		return
	}

	refreshTokenData, err := h.services.Authorization.ParseRefreshToken(refreshToken)

	if err != nil {
		h.setRefreshToken(c, "", -1)
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.ID)
	_, err = h.services.Redis.GetRefreshToken(c.Request.Context(), key)

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "token does not exist")
		return
	}

	h.createTokens(c, refreshTokenData.Username, refreshTokenData.UserId)
}

func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cannot logout: invalid refreshToken cookie")
		return
	}

	refreshTokenData, err := h.services.Authorization.ParseRefreshToken(refreshToken)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid refreshToken")
		return
	}

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.ID)
	err = h.services.Redis.DeleteRefreshToken(c.Request.Context(), key)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.setRefreshToken(c, "", -1)
	c.Status(http.StatusOK)
}

func (h *Handler) createTokens(c *gin.Context, username string, userId uint) {
	accessToken, err := h.services.Authorization.GenerateAccessToken(username, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenData, err := h.services.Authorization.GenerateRefreshToken(username, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	key := fmt.Sprintf("%s:%s", username, refreshTokenData.ID)
	err = h.services.Redis.SetRefreshToken(c.Request.Context(), key, refreshTokenData.RefreshToken)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.setRefreshToken(c, refreshTokenData.RefreshToken, int(service.RefreshTokenTTL))
	c.JSON(http.StatusOK, map[string]string{
		"accessToken": accessToken,
	})
}

func (h *Handler) setRefreshToken(c *gin.Context, refreshToken string, TTL int) {
	c.SetCookie("refreshToken", refreshToken, TTL, "/api", "localhost", false, true)
}
