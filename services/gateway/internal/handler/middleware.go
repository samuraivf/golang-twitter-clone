package handler

import (
	"errors"
	"gateway/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authService "auth/proto"
)

const (
	authorizationHeader = "Authorization"
	userDataCtx         = "userData"
)

const (
	errUserIsAuthorized    = "user is authorized"
	errInvalidAuthHeader   = "invalid Authorization header"
	errTokenIsEmpty        = "token is empty"
	errUserNotFound        = "user not found"
	errInvalidUserDataType = "userData is of invalid type"
)

func (h *Handler) isUnauthorized(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	_, err := c.Cookie("refreshToken")

	if header != "" || err == nil {
		newErrorResponse(c, http.StatusBadRequest, errUserIsAuthorized)
		return
	}
}

func (h *Handler) isAuthorized(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, errInvalidAuthHeader)
		return
	}

	headerParts := strings.Split(header, " ")

	if headerParts[0] != "Bearer" || len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, errInvalidAuthHeader)
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, errTokenIsEmpty)
		return
	}

	connection := services.ConnectAuthGrpc()
	defer connection.Close()

	authClient := authService.NewAuthorizationClient(connection)

	tokenData, err := authClient.ParseAccessToken(c, &authService.AccessToken{
		AccessToken: headerParts[1],
	})
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userDataCtx, tokenData)
}

func getUserData(c *gin.Context) (*authService.TokenData, error) {
	userData, ok := c.Get(userDataCtx)

	if !ok {
		return nil, errors.New(errUserNotFound)
	}

	tokenData, ok := userData.(*authService.TokenData)

	if !ok {
		return nil, errors.New(errInvalidUserDataType)
	}

	return tokenData, nil
}
