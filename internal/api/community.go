package api

import (
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/service"
	"demoProject/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CommunityHandler struct {
	communityService service.ICommunityService
}

func NewcommunityHandler(communityService service.ICommunityService) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
	}
}

// CommunityListHandler
// @Summary 获取所有社区信息
// @Tags Community
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/communities [get]
func (comm *CommunityHandler) CommunityListHandler(c *gin.Context) {
	commListRes, err := comm.communityService.GetCommunityList(c)
	if err != nil {
		zap.L().Error("communityService.GetCommunityList failed", zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	c.JSON(http.StatusOK, util.JsonRsp(commListRes))

}

// CommunityDetailHandler
// @Summary 获取某个社区详细信息
// @Tags Community
// @Accept json
// @Produce json
// @Param id path int64 true "社区ID"
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/communities/{id} [get]
func (comm *CommunityHandler) CommunityDetailHandler(c *gin.Context) {
	var req dto.GetCommunityDetailReq

	if err := c.ShouldBindUri(&req); err != nil {
		zap.L().Error("c.ShouldBindQuery failed", zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	commDetailRes, err := comm.communityService.GetCommunityDetail(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	c.JSON(http.StatusOK, util.JsonRsp(commDetailRes))
}
