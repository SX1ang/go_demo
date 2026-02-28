package repo

import (
	"context"
	"demoProject/dal/model"
)

type IUserRepo interface {
	UserIsExist(ctx context.Context, username string) (bool, error)
	AddUser(ctx context.Context, user *model.User) error
}
