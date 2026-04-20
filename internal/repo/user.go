package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/infra"
	"errors"

	"gorm.io/gorm"
)

type MysqlUserRespository struct {
	db *gorm.DB
}

func NewMysqlUserRespository(res *infra.Resources) IUserRepo {
	return &MysqlUserRespository{
		db: res.DB,
	}
}

// UserIsExist 检查用户是否存在
func (m *MysqlUserRespository) UserIsExist(ctx context.Context, username string) (bool, error) {
	q := query.Use(m.db)

	// 用表对象的 WithContext
	_, err := q.User.
		WithContext(ctx).
		Where(q.User.Username.Eq(username)).
		First()

	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// AddUser 新增User记录
func (m *MysqlUserRespository) AddUser(ctx context.Context, user *model.User) error {
	err := query.Use(m.db).WithContext(ctx).User.Create(user)

	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlUserRespository) GetUser(ctx context.Context, username string) (*model.User, error) {
	q := query.Use(m.db)

	// 用表对象的 WithContext
	user, err := q.User.
		WithContext(ctx).
		Where(q.User.Username.Eq(username)).
		First()

	if err == nil {
		return user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}
