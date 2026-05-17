package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/infra"
	"demoProject/pkg/util"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
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

type PostRepository struct {
	db     *gorm.DB
	client *redis.Client
}

func NewPostRepository(res *infra.Resources) IPostRepo {
	return &PostRepository{
		db:     res.DB,
		client: res.Redis,
	}
}

func (m *PostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	err := query.Use(m.db).WithContext(ctx).Post.Create(post)
	if err != nil {
		return err
	}

	pipeline := m.client.TxPipeline()
	postIdStr := strconv.FormatInt(post.PostID, 10)

	pipeline.ZAdd(ctx, KeyPostTime, redis.Z{
		Score:  float64(post.CreateTime.Unix()),
		Member: postIdStr,
	})
	pipeline.ZAdd(ctx, KeyPostScore, redis.Z{
		Score:  0,
		Member: postIdStr,
	})

	_, err = pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostRepository) GetPostDetail(ctx context.Context, postId int64) (*dto.PostDetail, error) {

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

func (m *PostRepository) GetPostList(ctx context.Context, ids []int64) ([]*dto.PostDetail, error) {
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
		Where(p.PostID.In(ids...)).
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

func (m *PostRepository) PostIsExist(ctx context.Context, postId int64) (bool, error) {
	q := query.Use(m.db)

	_, err := q.WithContext(ctx).Post.Where(q.Post.PostID.Eq(postId)).First()

	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func (m *PostRepository) GetPostIdsInOrder(ctx context.Context, page, size int, order string) ([]string, error) {
	start := int64((page - 1) * size)
	stop := start + int64(size-1)

	// 从redis中获取所有帖子的id
	if order == "new" {
		// 按照时间排序
		return m.client.ZRevRange(ctx, KeyPostTime, start, stop).Result()
	} else {
		// 按照热度分数排序
		return m.client.ZRevRange(ctx, KeyPostScore, start, stop).Result()
	}
}
