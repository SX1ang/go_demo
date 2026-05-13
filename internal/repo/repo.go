package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/internal/api/dto"
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

type ICommunityRepo interface {
	GetCommunityList(ctx context.Context) ([]*CommunityBasic, error)
	GetCommunityDetail(ctx context.Context, communityId int64) (*model.Community, error)
	AddCommunity(ctx context.Context, community *model.Community) error
	CommunityIsExist(ctx context.Context, communityId int32) (bool, error)
}

type IPostRepo interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostDetail(ctx context.Context, postId int64) (*dto.PostDetail, error)
}
