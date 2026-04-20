package repo

import (
	"context"
	"demoProject/dal/model"
)

type IUserRepo interface {
	UserIsExist(ctx context.Context, username string) (bool, error)
	AddUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, username string) (*model.User, error)
}

type ISessionRepo interface {
	SaveSession(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, sessionId string) (*model.Session, error)
	DeleteSession(ctx context.Context, userId int64) error
	RevokeRefreshToken(ctx context.Context, userId int64) error
}
