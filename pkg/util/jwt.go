package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 过期时间
const RefreshTokenExpireDuration = time.Hour * 24
const AccessTokenExpireDuration = time.Minute * 10

// var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var jwtSecret = []byte("go-demo-secret")

// ParseWithClaims参数keyFunc
var keyFunc = func(token *jwt.Token) (interface{}, error) {
	// 检查算法
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
	}

	// 返回密钥给ParseWithClaims解析token使用
	return jwtSecret, nil
}

// Token携带的声明信息
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
// 生成access token和refresh token
func GenerateToken(userID int64) (accesssToken, refreshToken string, err error) {
	// 生成access token
	claims := Claims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(), // 过期时间
			Issuer:    "go-demo",
		},
	}

	accesssToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// refresh token不需要任何自定义数据
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    "go-demo",
	}).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accesssToken, refreshToken, nil
}

// VerifyToken verify token
// 验证access token
func VerifyToken(tokenStr string) (*Claims, error) {
	// 自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// 验证token有效性并断言Claims类型
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}

func RefreshToken(accesssToken, refreshToken string) (newAccesssToken, newRefreshToken string, err error) {
	// refresh token无效直接返回
	if _, err := jwt.Parse(refreshToken, keyFunc); err != nil {
		return "", "", err
	}

	var claims Claims
	_, err = jwt.ParseWithClaims(accesssToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GenerateToken(claims.UserID)
	}

	return accesssToken, refreshToken, nil
}
