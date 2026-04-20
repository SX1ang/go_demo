package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{
		secretKey: secretKey,
	}
}

func (maker *JWTMaker) GenerateToken(userId int64, tokenType string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := newUserClaims(userId, tokenType, duration)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate user claims")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate signed token")
	}

	return signedToken, claims, nil
}

func (maker *JWTMaker) VerifyToken(signedToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// verify token singing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
