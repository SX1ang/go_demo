package api

import (
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/service"
	"demoProject/pkg/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VoteHandler struct {
	voteService service.IVoteService
}

func NewVoteHandler(voteService service.IVoteService) *VoteHandler {
	return &VoteHandler{
		voteService: voteService,
	}
}

// VotePostHandler
// @Summary 给帖子投票
// @Tags Vote
// @Accept json
// @Produce json
// @Param data body dto.VotePostReq true "投票信息"
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/votes [post]
func (v *VoteHandler) VotePostHandler(c *gin.Context) {
	var req dto.VotePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("VotePostHandler ShouldBindJSON error", zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	err := v.voteService.VotePost(c, &req)
	if err != nil {
		var customErr *util.CustomizedErr
		if errors.As(err, &customErr) {
			c.JSON(http.StatusOK, util.JsonRsp(customErr))
			return
		}
		
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	c.JSON(http.StatusOK, util.JsonRsp("ok"))

}
