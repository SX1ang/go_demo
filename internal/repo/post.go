package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/infra"
	"demoProject/pkg/util"
	"time"

	"gorm.io/gorm"
)

type postDetailFlat struct {
	PostID         int64     `gorm:"column:post_id"`
	Title          string    `gorm:"column:title"`
	Content        string    `gorm:"column:content"`
	Status         int32     `gorm:"column:status"`
	PostCreateTime time.Time `gorm:"column:post_create_time"`
	PostUpdateTime time.Time `gorm:"column:post_update_time"`
	AuthorName     string    `gorm:"column:author_name"`
	CommunityID    int64     `gorm:"column:community_id"`
	CommunityName  string    `gorm:"column:community_name"`
	Introduction   string    `gorm:"column:introduction"`
	CommCreateTime time.Time `gorm:"column:comm_create_time"`
	CommUpdateTime time.Time `gorm:"column:comm_update_time"`
}

type MysqlPostRepository struct {
	db *gorm.DB
}

func NewMysqlPostRepository(res *infra.Resources) IPostRepo {
	return &MysqlPostRepository{db: res.DB}
}

func (m *MysqlPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	err := query.Use(m.db).WithContext(ctx).Post.Create(post)
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlPostRepository) GetPostDetail(ctx context.Context, postId int64) (*dto.PostDetail, error) {

	q := query.Use(m.db)
	p := q.Post
	u := q.User
	c := q.Community

	var flat postDetailFlat
	err := p.WithContext(ctx).
		Select(
			p.PostID.As("post_id"),
			p.Title,
			p.Content,
			p.Status,
			p.CreateTime.As("post_create_time"),
			p.UpdateTime.As("post_update_time"),
			u.Username.As("author_name"),
			c.CommunityID.As("community_id"),
			c.CommunityName.As("community_name"),
			c.Introduction,
			c.CreateTime.As("comm_create_time"),
			c.UpdateTime.As("comm_update_time"),
		).
		LeftJoin(u, p.AuthorID.EqCol(u.UserID)).
		LeftJoin(c, p.CommunityID.EqCol(c.CommunityID)).
		Where(p.PostID.Eq(postId)).
		Scan(&flat)

	// Scan查不到不报错，需要手动判断
	if flat.PostID == 0 {
		return nil, &util.CustomizedErr{
			Code: e.ERROR_NOT_EXIST_POST,
			Msg:  e.GetMsg(e.ERROR_NOT_EXIST_POST),
		}
	}

	if err != nil {
		return nil, err
	}

	return &dto.PostDetail{
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
	}, nil
}

func (m *MysqlPostRepository) GetPostList(ctx context.Context, page, size int) ([]*dto.PostDetail, error) {
	q := query.Use(m.db)
	p := q.Post
	u := q.User
	c := q.Community

	var flats []postDetailFlat
	err := p.WithContext(ctx).
		Select(
			p.PostID.As("post_id"),
			p.Title,
			p.Content,
			p.Status,
			p.CreateTime.As("post_create_time"),
			p.UpdateTime.As("post_update_time"),
			u.Username.As("author_name"),
			c.CommunityID.As("community_id"),
			c.CommunityName.As("community_name"),
			c.Introduction,
			c.CreateTime.As("comm_create_time"),
			c.UpdateTime.As("comm_update_time"),
		).
		LeftJoin(u, p.AuthorID.EqCol(u.UserID)).
		LeftJoin(c, p.CommunityID.EqCol(c.CommunityID)).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&flats)
	if err != nil {
		return nil, err
	}

	posts := make([]*dto.PostDetail, 0, len(flats))
	for _, flat := range flats {
		posts = append(posts, &dto.PostDetail{
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
		})
	}
	return posts, nil
}
