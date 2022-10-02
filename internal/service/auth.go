package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	day             = time.Hour * 24
	AccessTokenTTL  = day
	RefreshTokenTTL = day * 30
	jwtAccessKey    = "dghfugwep;fljafgfu29"
	jwtRefreshKey   = "hdfsgdfsjdfww0w0wvvd"
)

type TokenData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	UserId   uint   `json:"userId"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenData
}

type RefreshTokenData struct {
	RefreshToken string
	ID           string
}

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) GenerateAccessToken(username string, userId uint) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: username,
			UserId:   userId,
		},
	})

	return accessToken.SignedString([]byte(jwtAccessKey))
}

func (s *AuthService) GenerateRefreshToken(username string, userId uint) (*RefreshTokenData, error) {
	tokenID := uuid.NewString()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			ID:       tokenID,
			Username: username,
			UserId:   userId,
		},
	})

	token, err := refreshToken.SignedString([]byte(jwtRefreshKey))

	if err != nil {
		return nil, err
	}

	return &RefreshTokenData{RefreshToken: token, ID: tokenID}, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(jwtAccessKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return &claims.TokenData, nil
}

func (s *AuthService) ParseRefreshToken(refreshToken string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(jwtRefreshKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return &claims.TokenData, nil
}
