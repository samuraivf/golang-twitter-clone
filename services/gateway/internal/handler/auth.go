package handler

import (
	"fmt"
	"gateway/internal/dto"
	"gateway/internal/services"
	"net/http"
	"time"

	authService "auth/proto"
	redisService "redis/proto"
	userService "user/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
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

	userClient, closeUser := services.GetUserClient()
	defer closeUser()

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if _, err := userClient.GetUserByEmail(c, &userService.Email{Email: user.Email}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, errUserHaveAlreadyExists)
		return
	}

	if _, err := userClient.GetUserByUsername(c, &userService.Username{Username: user.Username}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, errUserHaveAlreadyExists)
		return
	}

	id, err := userClient.CreateUser(c, &userService.CreateUserDto{
		Username: user.Username,
		Password: user.Password,
		Name: user.Name,
		Email: user.Email,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]uint{
		"id": uint(id.UserId),
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

	userClient, closeUser := services.GetUserClient()
	defer closeUser()

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	userData, err := userClient.ValidateUser(c, &userService.ValidateUserParams{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	createTokens(c, userData.Username, uint(userData.Id))
}

func (h *Handler) refresh(c *gin.Context, createTokens createTokens) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		newErrorResponse(c, http.StatusBadRequest, errInvalidRefreshToken)
		return
	}

	authClient, closeAuth := services.GetAuthClient()
	defer closeAuth()

	refreshTokenData, err := authClient.ParseRefreshToken(c, &authService.RefreshToken{
		RefreshToken: refreshToken,
	})
	if err != nil {
		h.setRefreshToken(c, "", -1)
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	redisClient, closeRedis := services.GetRedisClient()
	defer closeRedis()

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.Id)
	_, err = redisClient.GetRefreshToken(c.Request.Context(), &redisService.Key{Key: key})
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, errTokenDoesNotExist)
		return
	}

	createTokens(c, refreshTokenData.Username, uint(refreshTokenData.UserId))
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

	authClient, closeAuth := services.GetAuthClient()
	defer closeAuth()

	refreshTokenData, err := authClient.ParseRefreshToken(c, &authService.RefreshToken{
		RefreshToken: refreshToken,
	})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidRefreshToken)
		return
	}

	redisClient, closeRedis := services.GetRedisClient()
	defer closeRedis()

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.Id)
	_, err = redisClient.DeleteRefreshToken(c.Request.Context(), &redisService.Key{Key: key})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.setRefreshToken(c, "", -1)
	c.Status(http.StatusOK)
}

func (h *Handler) createTokens(c *gin.Context, username string, userId uint) {
	authClient, closeAuth := services.GetAuthClient()
	defer closeAuth()

	accessToken, err := authClient.GenerateAccessToken(c, &authService.UserData{
		Username: username,
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenData, err := authClient.GenerateRefreshToken(c, &authService.UserData{
		Username: username,
		UserId: uint64(userId),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	redisClient, closeRedis := services.GetRedisClient()
	defer closeRedis()

	refreshTokenTTL, _ := authClient.GetRefreshTokenTTL(c, new(emptypb.Empty))

	key := fmt.Sprintf("%s:%s", username, refreshTokenData.Id)
	_, err = redisClient.SetRefreshToken(c, &redisService.SetRefreshTokenParams{
		Key: key,
		RefreshToken: refreshTokenData.RefreshToken,
		RefreshTokenTTL: durationpb.New(time.Duration(refreshTokenTTL.RefreshTokenTTL)),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.setRefreshToken(c, refreshTokenData.RefreshToken, int(refreshTokenTTL.RefreshTokenTTL))
	c.JSON(http.StatusOK, map[string]string{
		"accessToken": accessToken.AccessToken,
	})
}

func (h *Handler) setRefreshToken(c *gin.Context, refreshToken string, TTL int) {
	c.SetCookie("refreshToken", refreshToken, TTL, "/api", "localhost", false, true)
}
