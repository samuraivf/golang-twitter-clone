package service

import (
	"context"
	"errors"
	"time"

	pb "auth/proto"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	day             = time.Hour * 24
	AccessTokenTTL  = day
	refreshTokenTTL = day * 30
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
	pb.UnimplementedAuthorizationServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) GenerateAccessToken(context context.Context, userData *pb.UserData) (*pb.AccessToken, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: userData.Username,
			UserId:   uint(userData.UserId),
		},
	})

	token, err := accessToken.SignedString([]byte(jwtAccessKey))

	if err != nil {
		return new(pb.AccessToken), err
	}

	return &pb.AccessToken{
		AccessToken: token,
	}, nil
}

func (s *AuthService) GenerateRefreshToken(context context.Context, userData *pb.UserData) (*pb.RefreshTokenData, error) {
	tokenID := uuid.NewString()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			ID:       tokenID,
			Username: userData.Username,
			UserId:   uint(userData.UserId),
		},
	})

	token, err := refreshToken.SignedString([]byte(jwtRefreshKey))

	if err != nil {
		return new(pb.RefreshTokenData), err
	}

	return &pb.RefreshTokenData{RefreshToken: token, Id: tokenID}, nil
}

func (s *AuthService) ParseAccessToken(context context.Context, accessToken *pb.AccessToken) (*pb.TokenData, error) {
	token, err := jwt.ParseWithClaims(accessToken.AccessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return new(pb.TokenData), errors.New("invalid signing method")
		}

		return []byte(jwtAccessKey), nil
	})

	if err != nil {
		return new(pb.TokenData), err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return new(pb.TokenData), errors.New("token claims are not of type *tokenClaims")
	}

	return &pb.TokenData{
		UserId: uint64(claims.TokenData.UserId),
		Username: claims.TokenData.Username,
		Id: claims.TokenData.ID,
	}, nil
}

func (s *AuthService) ParseRefreshToken(context context.Context, refreshToken *pb.RefreshToken) (*pb.TokenData, error) {
	token, err := jwt.ParseWithClaims(refreshToken.RefreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return new(pb.TokenData), errors.New("invalid signing method")
		}

		return []byte(jwtRefreshKey), nil
	})

	if err != nil {
		return new(pb.TokenData), err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return new(pb.TokenData), errors.New("token claims are not of type *tokenClaims")
	}

	return &pb.TokenData{
		UserId: uint64(claims.TokenData.UserId),
		Username: claims.TokenData.Username,
		Id: claims.TokenData.ID,
	}, nil
}

func (s *AuthService) GetRefreshTokenTTL(context context.Context, _ *emptypb.Empty) (*pb.RefreshTokenTTL, error) {
	return &pb.RefreshTokenTTL{RefreshTokenTTL: uint64(refreshTokenTTL)}, nil
}
