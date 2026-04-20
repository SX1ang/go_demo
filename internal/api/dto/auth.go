package dto

import "time"

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type RevokeRefreshTokenReq struct {
	UserId int64
}
