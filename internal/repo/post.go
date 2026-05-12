package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"gorm.io/gorm"
)

type MysqlPostRepository struct {
	db *gorm.DB
}

func NewMysqlPostRepository(db *gorm.DB) IPostRepo {
	return &MysqlPostRepository{db: db}
}

func (m *MysqlPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	err := query.Use(m.db).WithContext(ctx).Post.Create(post)
	if err != nil {
		return err
	}
	return nil
}
