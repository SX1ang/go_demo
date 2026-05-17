package service

import (
	"demoProject/dal/model"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/repo"
	"demoProject/pkg/snowflake"
	"demoProject/pkg/util"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 在 model 或者一个单独的 constants 包里定义
const (
	PostStatusPending   int32 = iota + 1 // 审核中
	PostStatusPublished                  // 已发布
)

func PostStatusToText(s int32) string {
	switch s {
	case PostStatusPending:
		return "审核中"
	case PostStatusPublished:
		return "已发布"
	default:
		return "未知状态"
	}
}

type PostService struct {
	postRepo repo.IPostRepo
	commRepo repo.ICommunityRepo
	voteRepo repo.IVoteRepo
}

func NewPostService(postRepo repo.IPostRepo, commRepo repo.ICommunityRepo, voteRepo repo.IVoteRepo) IPostService {
	return &PostService{
		postRepo: postRepo,
		commRepo: commRepo,
		voteRepo: voteRepo,
	}
}

func (p *PostService) CreatePost(c *gin.Context, req *dto.CreatePostReq) error {
	// 1.检查community ID是否存在
	isExist, err := p.commRepo.CommunityIsExist(c, int32(req.CommunityID))
	if err != nil {
		zap.L().Error("check CommunityIsExist error", zap.Error(err))
		return err
	}

	if !isExist {
		zap.L().Error("CommunityIsExist false")
		return &util.CustomizedErr{
			Code: e.ERROR_NOT_EXIST_COMMUNITY,
			Msg:  e.GetMsg(e.ERROR_NOT_EXIST_COMMUNITY),
		}
	}

	// 2.构造post model
	postId := snowflake.GenID()
	authorId, exists := c.Get("UserId")
	if !exists {
		zap.L().Error("UserId not exist")
		return &util.CustomizedErr{
			Code: e.UNAUTHORIZED,
			Msg:  e.GetMsg(e.UNAUTHORIZED),
		}
	}

	post := model.Post{
		PostID:      postId,
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    authorId.(int64),
		CommunityID: req.CommunityID,
		Status:      PostStatusPending, // 默认审核中
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	err = p.postRepo.CreatePost(c, &post)
	if err != nil {
		zap.L().Error("postRepo.CreatePost error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostService) GetPostDetail(c *gin.Context, req *dto.GetPostDetailReq) (*dto.GetPostDetailRes, error) {
	postId := req.PostID
	postDetail, err := p.postRepo.GetPostDetail(c, postId)
	if err != nil {
		zap.L().Error("postRepo.GetPostDetail error", zap.Error(err))
		return nil, err
	}

	// 补充帖子点赞数量信息
	counts, err := p.voteRepo.GetPostsVoteCount(c, []int64{postId})
	if err != nil {
		zap.L().Error("voteRepo.GetPostsVoteCount error", zap.Error(err))
		return nil, err
	}
	postDetail.VoteCount = counts[0]

	postDetail.StatusText = PostStatusToText(postDetail.Status)

	return &dto.GetPostDetailRes{
		Post: postDetail,
	}, nil
}

func (p *PostService) GetPostList(c *gin.Context, req *dto.GetPostListReq) (*dto.GetPostListRes, error) {
	page := req.Page
	size := req.Size
	order := req.Order

	// 按照page, size, order获取所有帖子的id
	ids, err := p.postRepo.GetPostIdsInOrder(c, page, size, order)
	if err != nil {
		zap.L().Error("postRepo.GetPostIdsInOrder error", zap.Error(err))
		return nil, err
	}

	// 把 []string 转成 []int64
	intIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid id %q: %w", s, err)
		}
		intIDs = append(intIDs, v)
	}

	// 根据帖子ids获取帖子信息
	postList, err := p.postRepo.GetPostList(c, intIDs)
	if err != nil {
		zap.L().Error("postRepo.GetPostList error", zap.Error(err))
		return nil, err
	}

	// 补充帖子点赞数量信息
	counts, err := p.voteRepo.GetPostsVoteCount(c, intIDs)
	if err != nil {
		zap.L().Error("voteRepo.GetPostsVoteCount error", zap.Error(err))
		return nil, err
	}

	for idx, post := range postList {
		post.VoteCount = counts[idx]
	}

	return &dto.GetPostListRes{
		Posts: postList,
	}, nil
}
