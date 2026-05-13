package service

import (
	"context"
	"demoProject/internal/api/dto"

	"github.com/gin-gonic/gin"
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

type ICommunityService interface {
	GetCommunityList(c *gin.Context) (*dto.GetCommunityListRes, error)
	GetCommunityDetail(c *gin.Context, req *dto.GetCommunityDetailReq) (*dto.GetCommunityDetailRes, error)
}

type IPostService interface {
	CreatePost(c *gin.Context, req *dto.CreatePostReq) error
	GetPostDetail(c *gin.Context, req *dto.GetPostDetailReq) (*dto.GetPostDetailRes, error)
}
