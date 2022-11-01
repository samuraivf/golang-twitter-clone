package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/samuraivf/twitter-clone/internal/dto"
	"github.com/samuraivf/twitter-clone/internal/repository/models"
	"github.com/samuraivf/twitter-clone/internal/service"
	mock_service "github.com/samuraivf/twitter-clone/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehaviour func (s *mock_service.MockUser, user dto.CreateUserDto)

	tests := []struct{
		name string
		inputBody string
		inputUser dto.CreateUserDto
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{"username": "username", "password": "username", "email": "email@gmail.com", "name": "Username"}`,
			inputUser: dto.CreateUserDto{
				LoginDto: dto.LoginDto{
					Username: "username",
					Password: "username",
				},
				Email: "email@gmail.com",
				Name: "Username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.CreateUserDto) {
				s.EXPECT().GetUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound)
				s.EXPECT().GetUserByUsername(user.Username).Return(nil, gorm.ErrRecordNotFound)
				s.EXPECT().CreateUser(user).Return(uint(1), nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name: "Empty Fields",
			inputBody: `{"email": "email@gmail.com", "name": "Username"}`,
			inputUser: dto.CreateUserDto{},
			mockBehaviour: func(s *mock_service.MockUser, user dto.CreateUserDto) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"error invalid input body"}`,
		},
		{
			name: "Service fail",
			inputBody: `{"username": "username", "password": "username", "email": "email@gmail.com", "name": "Username"}`,
			inputUser: dto.CreateUserDto{
				LoginDto: dto.LoginDto{
					Username: "username",
					Password: "username",
				},
				Email: "email@gmail.com",
				Name: "Username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.CreateUserDto) {
				s.EXPECT().GetUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound)
				s.EXPECT().GetUserByUsername(user.Username).Return(nil, gorm.ErrRecordNotFound)
				s.EXPECT().CreateUser(user).Return(uint(1), errors.New("service fail"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service fail"}`,
		},
		{
			name: "GetUserByEmail User exists",
			inputBody: `{"username": "username", "password": "username", "email": "email@gmail.com", "name": "Username"}`,
			inputUser: dto.CreateUserDto{
				LoginDto: dto.LoginDto{
					Username: "username",
					Password: "username",
				},
				Email: "email@gmail.com",
				Name: "Username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.CreateUserDto) {
				s.EXPECT().GetUserByEmail(user.Email).Return(&models.User{}, nil)
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"user have already exists"}`,
		},
		{
			name: "GetUserByUsername UserExists",
			inputBody: `{"username": "username", "password": "username", "email": "email@gmail.com", "name": "Username"}`,
			inputUser: dto.CreateUserDto{
				LoginDto: dto.LoginDto{
					Username: "username",
					Password: "username",
				},
				Email: "email@gmail.com",
				Name: "Username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.CreateUserDto) {
				s.EXPECT().GetUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound)
				s.EXPECT().GetUserByUsername(user.Username).Return(&models.User{}, nil)
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"user have already exists"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehaviour func (s *mock_service.MockUser, user dto.LoginDto)

	tests := []struct{
		name string
		inputBody string
		inputUser dto.LoginDto
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{"username": "username", "password": "username"}`,
			inputUser: dto.LoginDto{
				Username: "username",
				Password: "username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.LoginDto) {
				s.EXPECT().ValidateUser(user.Username, user.Password).Return(&models.User{
					Model: models.Model{
						ID: 1,
					},
					Username: "username",
					Password: "username",
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"accessToken":"access_token"}`,
		},
		{
			name: "Error in Validate User",
			inputBody: `{"username": "username", "password": "username"}`,
			inputUser: dto.LoginDto{
				Username: "username",
				Password: "username",
			},
			mockBehaviour: func (s *mock_service.MockUser, user dto.LoginDto) {
				s.EXPECT().ValidateUser(user.Username, user.Password).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"some error"}`,
		},
		{
			name: "Invalid input body",
			inputBody: `{"username": "username"}`,
			mockBehaviour: func (s *mock_service.MockUser, user dto.LoginDto) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"error invalid input body"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehaviour(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/sign-in", func(c *gin.Context) {
				handler.signIn(c, func(c *gin.Context, username string, userId uint) {
					c.JSON(http.StatusOK, map[string]string{
						"accessToken": "access_token",
					})
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_refresh(t *testing.T) {
	type mockBehaviour func (a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string)

	tests := []struct{
		name string
		refreshTokenCookie string
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(&service.TokenData{
					ID: "asd",
					Username: "username",
					UserId: 1,
				}, nil)

				ctx := context.Background()
				r.EXPECT().GetRefreshToken(ctx, "username:asd").Return("asd", nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"accessToken":"access_token"}`,
		},
		{
			name: "Invalid refresh token",
			refreshTokenCookie: "",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid refresh token"}`,
		},
		{
			name: "Error in ParseRefreshToken",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"some error"}`,
		},
		{
			name: "Error in GetRefreshToken",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(&service.TokenData{
					ID: "asd",
					Username: "username",
					UserId: 1,
				}, nil)

				ctx := context.Background()
				r.EXPECT().GetRefreshToken(ctx, "username:asd").Return("", errors.New("some error"))
			},
			expectedStatusCode: 401,
			expectedResponseBody: `{"message":"token does not exist"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			redis := mock_service.NewMockRedis(c)

			test.mockBehaviour(auth, redis, test.refreshTokenCookie)

			services := &service.Service{Authorization: auth, Redis: redis}
			handler := Handler{services}

			r := gin.New()
			r.GET("/refresh", func(c *gin.Context) {
				handler.refresh(c, func(c *gin.Context, username string, userId uint) {
					c.JSON(http.StatusOK, map[string]string{
						"accessToken": "access_token",
					})
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/refresh", nil)

			cookie := http.Cookie{
				Name: "refreshToken",
				Value: test.refreshTokenCookie,
			}
			req.AddCookie(&cookie)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_logout(t *testing.T) {
	type mockBehaviour func (a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string)

	tests := []struct{
		name string
		refreshTokenCookie string
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(&service.TokenData{
					ID: "asd",
					Username: "username",
					UserId: 1,
				}, nil)

				ctx := context.Background()
				r.EXPECT().DeleteRefreshToken(ctx, "username:asd").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: "",
		},
		{
			name: "Invalid refresh token",
			refreshTokenCookie: "",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid refresh token"}`,
		},
		{
			name: "Error in ParseRefreshToken",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid refresh token"}`,
		},
		{
			name: "Error in DeleteRefreshToken",
			refreshTokenCookie: "refresh_token",
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, token string) {
				a.EXPECT().ParseRefreshToken(token).Return(&service.TokenData{
					ID: "asd",
					Username: "username",
					UserId: 1,
				}, nil)

				ctx := context.Background()
				r.EXPECT().DeleteRefreshToken(ctx, "username:asd").Return(errors.New("some error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			redis := mock_service.NewMockRedis(c)

			test.mockBehaviour(auth, redis, test.refreshTokenCookie)

			services := &service.Service{Authorization: auth, Redis: redis}
			handler := Handler{services}

			r := gin.New()
			r.GET("/logout", handler.logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)

			cookie := http.Cookie{
				Name: "refreshToken",
				Value: test.refreshTokenCookie,
			}
			req.AddCookie(&cookie)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_createTokens(t *testing.T) {
	type mockBehaviour func (a *mock_service.MockAuthorization, r *mock_service.MockRedis, username string, userId uint)

	tests := []struct{
		name string
		username string
		userId uint
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "OK",
			username: "username",
			userId: uint(1),
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, username string, userId uint) {
				a.EXPECT().GenerateAccessToken(username, userId).Return("access_token", nil)
				a.EXPECT().GenerateRefreshToken(username, userId).Return(&service.RefreshTokenData{
					RefreshToken: "refresh_token",
					ID: "dshg27",
				}, nil)
				ctx := context.Background()
				r.EXPECT().SetRefreshToken(ctx, "username:dshg27", "refresh_token").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"accessToken":"access_token"}`,
		},
		{
			name: "Error in GenerateAccessToken",
			username: "username",
			userId: uint(1),
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, username string, userId uint) {
				a.EXPECT().GenerateAccessToken(username, userId).Return("", errors.New("some error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"some error"}`,
		},
		{
			name: "Error in GenerateRefreshToken",
			username: "username",
			userId: uint(1),
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, username string, userId uint) {
				a.EXPECT().GenerateAccessToken(username, userId).Return("access_token", nil)
				a.EXPECT().GenerateRefreshToken(username, userId).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"some error"}`,
		},
		{
			name: "Error in SetRefreshToken",
			username: "username",
			userId: uint(1),
			mockBehaviour: func(a *mock_service.MockAuthorization, r *mock_service.MockRedis, username string, userId uint) {
				a.EXPECT().GenerateAccessToken(username, userId).Return("access_token", nil)
				a.EXPECT().GenerateRefreshToken(username, userId).Return(&service.RefreshTokenData{
					RefreshToken: "refresh_token",
					ID: "dshg27",
				}, nil)
				ctx := context.Background()
				r.EXPECT().SetRefreshToken(ctx, "username:dshg27", "refresh_token").Return(errors.New("some error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", nil)

			r := gin.New()

			auth := mock_service.NewMockAuthorization(c)
			redis := mock_service.NewMockRedis(c)
			test.mockBehaviour(auth, redis, test.username, test.userId)

			services := &service.Service{Authorization: auth, Redis: redis}
			handler := Handler{services}

			r.POST("/", func(ctx *gin.Context) {
				handler.createTokens(ctx, test.username, test.userId)
			})

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_setRefreshToken(t *testing.T) {
	t.Run("Set refresh token to cookie", func(t *testing.T) {
		r := gin.New()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api", nil)

		handler := Handler{}

		r.GET("/api", func(ctx *gin.Context) {
			handler.setRefreshToken(ctx, "refresh_token", 10000)

			respWriter := ctx.Writer.Header().Values("Set-Cookie")
			cookieData := strings.Split(respWriter[0], "; ")

			var cookieMap = map[string]string{}

			for _, data := range cookieData {
				keyValuePair := strings.Split(data, "=")

				if len(keyValuePair) == 1 {
					cookieMap[keyValuePair[0]] = "true"
					continue
				}
				cookieMap[keyValuePair[0]] = keyValuePair[1]
			}

			assert.Equal(t, cookieMap["refreshToken"], "refresh_token")
			assert.Equal(t, cookieMap["Max-Age"], "10000")
			assert.Equal(t, cookieMap["Path"], "/api")
			assert.Equal(t, cookieMap["Domain"], "localhost")
			assert.Equal(t, cookieMap["HttpOnly"], "true")
		})
		r.ServeHTTP(w, req)
	})
}
