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

type PostHandler struct {
	postService service.IPostService
}

func NewPostHandler(postService service.IPostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePostHandler
// @Summary 发布帖子
// @Tags Post
// @Accept json
// @Produce json
// @Param data body dto.CreatePostReq true "创建的帖子信息"
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/posts [post]
func (p *PostHandler) CreatePostHandler(ctx *gin.Context) {
	// 1.获取请求参数
	var req = new(dto.CreatePostReq)
	if err := ctx.ShouldBindJSON(&req); err != nil { // ShouldBindJSON contain a validator -> struct binding tag
		zap.L().Error("ShouldBindJSON failed", zap.Error(err))
		ctx.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))

		return
	}

	// 2.创建帖子
	err := p.postService.CreatePost(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))

		return
	}

	// 3.返回响应
	ctx.JSON(http.StatusOK, util.JsonRsp("ok"))

}
