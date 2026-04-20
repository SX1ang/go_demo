package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	TokenTypeAccessToken  = "access_token"
	TokenTypeRefreshToken = "refresh_token"
)

type UserClaims struct {
	UserId    int64  `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

func newUserClaims(userId int64, tokenType string, duration time.Duration) (*UserClaims, error) {
	refreshTokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token id: %v", err)
	}

	return &UserClaims{
		UserId:    userId,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			Id:        refreshTokenId.String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}, nil
}
