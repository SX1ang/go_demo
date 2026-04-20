package service

import (
	"context"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/repo"
	"demoProject/pkg/util"
	"time"

	"go.uber.org/zap"
)

type AuthService struct {
	sessionRepo repo.ISessionRepo
	tokenMaker  *util.JWTMaker
}

func NewAuthService(sessionRepo repo.ISessionRepo, tokenMaker *util.JWTMaker) IAuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		tokenMaker:  tokenMaker,
	}
}

func (a *AuthService) RenewAccessToken(ctx context.Context, req *dto.RenewAccessTokenReq) (*dto.RenewAccessTokenRes, error) {
	refreshClaims, err := a.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		zap.L().Error("renew access token fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			Msg:  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
		}
	}

	// check token type, must refresh token
	if refreshClaims.TokenType != util.TokenTypeRefreshToken {
		zap.L().Error("renew access token fail, token type is not refresh token")
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	// get session info
	sessionId := refreshClaims.StandardClaims.Id
	session, err := a.sessionRepo.GetSession(ctx, sessionId)
	if err != nil || session == nil {
		zap.L().Error("renew access token fail, get session failed", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	// 检查refresh token是否被吊销
	if session.IsRevoked {
		zap.L().Error("renew access token fail, session is revoked")
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH_CHECK_TOKEN_EXPIRED,
			Msg:  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_EXPIRED),
		}
	}

	// 检查user id
	if session.UserID != refreshClaims.UserId {
		zap.L().Error("renew access token fail, user is wrong")
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	// create access token
	accessToken, accessClaims, err := a.tokenMaker.GenerateToken(session.UserID, util.TokenTypeAccessToken, 15*time.Minute)
	if err != nil {
		zap.L().Error("login step, generate token fail", zap.Error(err))
		return nil, &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}

	return &dto.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: time.Unix(accessClaims.StandardClaims.ExpiresAt, 0),
	}, nil
}

func (a *AuthService) RevokeRefreshToken(ctx context.Context, req *dto.RevokeRefreshTokenReq) error {
	err := a.sessionRepo.RevokeRefreshToken(ctx, req.UserId)
	if err != nil {
		zap.L().Error("revoke refresh token failed", zap.Error(err))
		return &util.CustomizedErr{
			Code: e.ERROR_AUTH,
			Msg:  e.GetMsg(e.ERROR_AUTH),
		}
	}
	return nil
}
