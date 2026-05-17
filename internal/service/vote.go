package service

import (
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/repo"
	"demoProject/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VoteService struct {
	voteRepo repo.IVoteRepo
	postRepo repo.IPostRepo
}

func NewVoteService(voteRepo repo.IVoteRepo, postRepo repo.IPostRepo) IVoteService {
	return &VoteService{
		voteRepo: voteRepo,
		postRepo: postRepo,
	}
}

func (v VoteService) VotePost(c *gin.Context, req *dto.VotePostReq) error {
	isExist, err := v.postRepo.PostIsExist(c, req.PostID)
	if err != nil {
		zap.L().Error("check PostIsExist error", zap.Error(err))
		return err
	}

	if !isExist {
		zap.L().Error("PostIsExist false")
		return &util.CustomizedErr{
			Code: e.ERROR_NOT_EXIST_POST,
			Msg:  e.GetMsg(e.ERROR_NOT_EXIST_POST),
		}
	}

	userId, exists := c.Get("UserId")
	if !exists {
		zap.L().Error("UserId not exist")
		return &util.CustomizedErr{
			Code: e.UNAUTHORIZED,
			Msg:  e.GetMsg(e.UNAUTHORIZED),
		}
	}

	err = v.voteRepo.VotePost(c, userId.(int64), req.PostID, req.Vote)
	if err != nil {
		zap.L().Error("voteRepo VotePost error", zap.Error(err))
		return err
	}

	return nil
}
