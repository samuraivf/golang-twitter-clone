package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/service"
)

const (
	authorizationHeader = "Authorization"
	userDataCtx         = "userData"
)

func (h *Handler) isUnauthorized(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	_, err := c.Cookie("refreshToken")

	if header != "" || err == nil {
		newErrorResponse(c, http.StatusBadRequest, "user is authorized")
		return
	}
}

func (h *Handler) isAuthorized(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if headerParts[0] != "Bearer" || len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	tokenData, err := h.services.Authorization.ParseAccessToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userDataCtx, tokenData)
}

func getUserData(c *gin.Context) (*service.TokenData, error) {
	userData, ok := c.Get(userDataCtx)

	if !ok {
		return nil, errors.New("user is not found")
	}

	tokenData, ok := userData.(*service.TokenData)

	if !ok {
		return nil, errors.New("userData is of invalid type")
	}

	return tokenData, nil
}
