package service

import (
	"context"
	"demoProject/dal/model"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/repo"
	"demoProject/pkg/util"
	"fmt"
	"time"

	"go.uber.org/zap"
)
import "demoProject/pkg/snowflake"

// 实现IUserService接口
type UserService struct {
	userRepo    repo.IUserRepo
	sessionRepo repo.ISessionRepo
	tokenMaker  *util.JWTMaker
}

// 接口类型本身就是一个“引用语义”的值，返回值不需要再取指针
func NewUserService(userRepo repo.IUserRepo, sessionRepo repo.ISessionRepo, tokenMaker *util.JWTMaker) IUserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenMaker:  tokenMaker,
	}

}

func (u *UserService) SignUp(ctx context.Context, req *dto.SignUpReq) error {
	// 1.判断用户存不存在
	is_exist, err := u.userRepo.UserIsExist(ctx, req.Username)
	if err != nil {
		return &util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	if is_exist {
		return &util.CustomizedErr{
			Code: e.ERROR_EXIST_USER,
			Msg:  e.GetMsg(e.ERROR_EXIST_USER),
		}
	}

	// 2.生成UUID
	uuid := snowflake.GenID()
	zap.L().Info(fmt.Sprintf("generate user ID:%s\n", uuid))

	// 3.保存用户信息至数据库
	err = u.userRepo.AddUser(ctx, &model.User{
		UserID:   uuid,
		Username: req.Username,
		Password: util.EncodeMD5(req.Password), // 敏感信息加密
		Email:    req.Email,
		Gender:   req.Gender,
	})

	if err != nil {
		zap.L().Error("add user fail", zap.Error(err))
		return &util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginRes, error) {
	// 验证用户名和密码
	user, err := u.userRepo.GetUser(ctx, req.Username)
	if err != nil {
		zap.L().Error("login step, get user fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	if user == nil {
		zap.L().Info("user not found", zap.String("username", req.Username))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_INVALID_CREDENTIALS,
			Msg:  e.GetMsg(e.ERROR_INVALID_CREDENTIALS),
		}
	}

	if user.Password != util.EncodeMD5(req.Password) {
		zap.L().Info("check password fail", zap.String("username", req.Username))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_INVALID_CREDENTIALS,
			Msg:  e.GetMsg(e.ERROR_INVALID_CREDENTIALS),
		}
	}

	// create json web token
	accessToken, accessClaims, err := u.tokenMaker.GenerateToken(user.UserID, util.TokenTypeAccessToken, 15*time.Minute)
	if err != nil {
		zap.L().Error("login step, generate token fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	refreshToken, refreshClaims, err := u.tokenMaker.GenerateToken(user.UserID, util.TokenTypeRefreshToken, 24*time.Hour)
	if err != nil {
		zap.L().Error("login step, generate token fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	session := &model.Session{
		SessionID:    refreshClaims.StandardClaims.Id,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Unix(refreshClaims.StandardClaims.ExpiresAt, 0),
	}

	err = u.sessionRepo.SaveSession(ctx, session)

	if err != nil {
		zap.L().Error("login step, create session fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	return &dto.LoginRes{
		UserID:        user.UserID,
		Username:      user.Username,
		SessionId:     session.SessionID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpire:  accessClaims.StandardClaims.ExpiresAt,
		RefreshExpire: session.ExpiresAt.Unix(),
	}, nil
}

func (u *UserService) Logout(ctx context.Context, req *dto.LogoutReq) error {
	userId := req.UserId
	err := u.sessionRepo.DeleteSession(ctx, userId)

	if err != nil {
		zap.L().Error("logout step, delete session fail", zap.Error(err))
		return &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	return nil
}
