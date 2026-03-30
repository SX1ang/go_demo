package service

import (
	"context"
	"demoProject/internal/api/dto"
)

type IUserService interface {
	SignUp(ctx context.Context, req *dto.SignUpReq) error
	Login(ctx context.Context, req *dto.LoginReq) error
}
