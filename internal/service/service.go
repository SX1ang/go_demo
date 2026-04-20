package service

import (
	"context"
	"demoProject/internal/api/dto"
)

type IUserService interface {
	SignUp(ctx context.Context, req *dto.SignUpReq) error
	Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginRes, error)
	Logout(ctx context.Context, req *dto.LogoutReq) error
}

type IAuthService interface {
	RenewAccessToken(ctx context.Context, req *dto.RenewAccessTokenReq) (*dto.RenewAccessTokenRes, error)
	RevokeRefreshToken(ctx context.Context, req *dto.RevokeRefreshTokenReq) error
}
