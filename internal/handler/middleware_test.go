package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/samuraivf/twitter-clone/internal/service"
	mock_service "github.com/samuraivf/twitter-clone/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_isAuthorized(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, token string)

	tests := []struct{
		name string
		headerName string
		headerValue string
		token string
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			headerName: "Authorization",
			headerValue: "Bearer token",
			token: "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(&service.TokenData{
					ID: "asdf",
					Username: "username",
					UserId: 1,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: "1",
		},
		{
			name: "Empty Authorization header",
			headerName: "",
			headerValue: "",
			token: "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"invalid Authorization header"}`,
		},
		{
			name: "Invalid Authorization header",
			headerName: "Authoriation",
			headerValue: "Bearer token",
			token: "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"invalid Authorization header"}`,
		},
		{
			name: "Invalid Bearer",
			headerName: "Authorization",
			headerValue: "Berer token",
			token: "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"invalid Authorization header"}`,
		},
		{
			name: "Empty token",
			headerName: "Authorization",
			headerValue: "Bearer ",
			token: "",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name: "Error in ParseAccessToken",
			headerName: "Authorization",
			headerValue: "Bearer token",
			token: "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			test.mockBehaviour(auth, test.token)

			services := &service.Service{Authorization: auth}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/auth", handler.isAuthorized, func(c *gin.Context) {
				userData, _ := getUserData(c)
				c.String(200, fmt.Sprintf("%d", userData.UserId))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/auth", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_isUnauthorized(t *testing.T) {
	tests := []struct{
		name string
		authHeader string
		authHeaderValue string
		refreshToken string
		refreshTokenValue string
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			expectedStatusCode: 200,
			expectedResponseBody: "",
		},
		{
			name: "Authorization header is not empty",
			authHeader: "Authorization",
			authHeaderValue: "Bearer dsfj",
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"user is authorized"}`,
		},
		{
			name: "Refresh token is not empty",
			refreshToken: "refreshToken",
			refreshTokenValue: "token",
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"user is authorized"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := Handler{}

			r := gin.New()
			r.GET("/unauth", handler.isUnauthorized, func(c *gin.Context) {
				c.String(200, "")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/unauth", nil)
			req.Header.Set(test.authHeader, test.authHeaderValue)
			cookie := http.Cookie{
				Name: test.refreshToken,
				Value: test.refreshTokenValue,
			}
			req.AddCookie(&cookie)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}