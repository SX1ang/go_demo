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
		var customErr *util.CustomizedErr
		if errors.As(err, &customErr) {
			ctx.JSON(http.StatusOK, util.JsonRsp(customErr))
		} else {
			ctx.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
				Code: e.ERROR,
				Msg:  e.GetMsg(e.ERROR),
			}))
		}

		return
	}

	// 3.返回响应
	ctx.JSON(http.StatusOK, util.JsonRsp("ok"))

}

// PostDetailHandler
// @Summary 获取某个帖子详细信息
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int64 true "帖子ID"
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/posts/{id} [get]
func (p *PostHandler) PostDetailHandler(ctx *gin.Context) {
	// 1.获取帖子ID
	var req = new(dto.GetPostDetailReq)
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		zap.L().Error("ShouldBindUri failed", zap.Error(err))
		ctx.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))
		return
	}

	// 2.根据帖子ID查询帖子详情
	postDetailRes, err := p.postService.GetPostDetail(ctx, req)
	if err != nil {
		var customErr *util.CustomizedErr
		if errors.As(err, &customErr) {
			ctx.JSON(http.StatusOK, util.JsonRsp(customErr))
			return
		}
		ctx.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	ctx.JSON(http.StatusOK, util.JsonRsp(postDetailRes))
}

// PostListHandler
// @Summary 分页获取所有帖子信息
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1) minimum(1)
// @Param size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Success 200 {object} util.JsonResult
// @Router /v1/posts [get]
func (p *PostHandler) PostListHandler(c *gin.Context) {
	var req dto.GetPostListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Error("ShouldBindQuery failed", zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))

		return
	}

	// 默认按照时间排序
	if req.Order == "" {
		req.Order = "new"
	}

	postList, err := p.postService.GetPostList(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.ERROR,
			Msg:  e.GetMsg(e.ERROR),
		}))
		return
	}

	c.JSON(http.StatusOK, util.JsonRsp(postList))
}
