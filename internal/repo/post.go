package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/api/dto"
	"time"

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

func (m *MysqlPostRepository) GetPostDetail(ctx context.Context, postId int64) (*dto.PostDetail, error) {
	var flat struct {
		PostID         int64     `db:"post_id"`
		Title          string    `db:"title"`
		Content        string    `db:"content"`
		Status         int32     `db:"status"`
		PostCreateTime time.Time `db:"post_create_time"`
		PostUpdateTime time.Time `db:"post_update_time"`
		AuthorName     string    `db:"author_name"`
		CommunityID    int64     `db:"community_id"`
		CommunityName  string    `db:"community_name"`
		Introduction   string    `db:"introduction"`
		CommCreateTime time.Time `db:"comm_create_time"`
		CommUpdateTime time.Time `db:"comm_update_time"`
	}

	err := m.db.WithContext(ctx).Raw(`SELECT
            p.post_id,
            p.title,
            p.content,
            p.status,
            p.create_time AS post_create_time,
            p.update_time AS post_update_time,
            u.username AS author_name,
            c.community_id AS community_id,
            c.community_name AS community_name,
            c.introduction AS introduction,
            c.create_time AS comm_create_time,
            c.update_time AS comm_update_time
        FROM post p
        LEFT JOIN user u ON p.author_id = u.user_id
        LEFT JOIN community c ON p.community_id = c.community_id
        WHERE p.post_id = ?
    `, postId).Scan(&flat).Error
	if err != nil {
		return nil, err
	}

	result := &dto.PostDetail{
		ID:         flat.PostID,
		Title:      flat.Title,
		Content:    flat.Content,
		Status:     flat.Status,
		CreateTime: flat.PostCreateTime,
		UpdateTime: flat.PostUpdateTime,
		AuthorName: flat.AuthorName,
		Community: dto.CommunityDetail{
			ID:           flat.CommunityID,
			Name:         flat.CommunityName,
			Introduction: flat.Introduction,
			CreateTime:   flat.CommCreateTime,
			UpdateTime:   flat.CommUpdateTime,
		},
	}

	return result, nil
}
