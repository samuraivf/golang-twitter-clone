package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/service"
)

type createTokens func(c *gin.Context, username string, userId uint)

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID sign-up
// @Accept json
// @Produce json
// @Param input body dto.CreateUserDto true "create user data"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var user dto.CreateUserDto

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if _, err := h.services.User.GetUserByEmail(user.Email); err == nil {
		newErrorResponse(c, http.StatusBadRequest, errUserHaveAlreadyExists)
		return
	}

	if _, err := h.services.User.GetUserByUsername(user.Username); err == nil {
		newErrorResponse(c, http.StatusBadRequest, errUserHaveAlreadyExists)
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

// @Summary SignIn
// @Tags auth
// @Description Sign in
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body dto.LoginDto true "login data"
// @Success 200 {string} string "token"
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context, createTokens createTokens) {
	var user dto.LoginDto

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	userData, err := h.services.User.ValidateUser(user.Username, user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	createTokens(c, userData.Username, userData.ID)
}

func (h *Handler) refresh(c *gin.Context, createTokens createTokens) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil || refreshToken == "" {
		newErrorResponse(c, http.StatusBadRequest, errInvalidRefreshToken)
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
		newErrorResponse(c, http.StatusUnauthorized, errTokenDoesNotExist)
		return
	}

	createTokens(c, refreshTokenData.Username, refreshTokenData.UserId)
}

// @Summary Logout
// @Tags auth
// @Description Logout
// @ID logout
// @Param cookie header string true "token"
// @Success 200
// @Failure 400 {object} ErrorMessage
// @Failure 500 {obkect} ErrorMessage
// @Router /auth/logout [get]
func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil || refreshToken == "" {
		newErrorResponse(c, http.StatusBadRequest, errInvalidRefreshToken)
		return
	}

	refreshTokenData, err := h.services.Authorization.ParseRefreshToken(refreshToken)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidRefreshToken)
		return
	}

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.ID)
	err = h.services.Redis.DeleteRefreshToken(c.Request.Context(), key)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
		return
	}

	h.setRefreshToken(c, refreshTokenData.RefreshToken, int(service.RefreshTokenTTL))
	c.JSON(http.StatusOK, map[string]string{
		"accessToken": accessToken,
	})
}

func (h *Handler) setRefreshToken(c *gin.Context, refreshToken string, TTL int) {
	c.SetCookie("refreshToken", refreshToken, TTL, "/api", "localhost", false, true)
}
