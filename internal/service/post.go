package service

import (
	"demoProject/dal/model"
	"demoProject/internal/api/dto"
	"demoProject/internal/repo"
	"demoProject/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// 在 model 或者一个单独的 constants 包里定义
const (
	PostStatusPending   int32 = iota // 审核中
	PostStatusPublished              // 已发布
)

type PostService struct {
	postRepo repo.IPostRepo
}

func NewPostService(postRepo repo.IPostRepo) IPostService {
	return &PostService{postRepo: postRepo}
}

func (p *PostService) CreatePost(c *gin.Context, req *dto.CreatePostReq) error {
	// 1.构造post model
	postId := snowflake.GenID()
	authorId := c.MustGet("UserId").(int64)

	post := model.Post{
		PostID:      postId,
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    authorId,
		CommunityID: req.CommunityID,
		Status:      PostStatusPending, // 默认审核中
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	err := p.postRepo.CreatePost(c, &post)
	if err != nil {
		zap.L().Error("postRepo.CreatePost error", zap.Error(err))
		return err
	}

	return nil
}
