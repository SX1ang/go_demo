package api

import (
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/service"
	"demoProject/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// SignUpHandler
// @Summary 用户注册
// @Tags User
// @Accept json
// @Produce json
// @Param data body dto.SignUpReq true "注册请求信息"
// @Success 200 {object} util.JsonResult
// @Router /v1/users/signup [post]
func (u *UserHandler) SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	var req = new(dto.SignUpReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("SignUpHandler", zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))
		return
	}

	// 2.业务处理
	// 请求级 context.Context，和每个请求关联
	ctx := c.Request.Context()
	// 也可以基于请求 ctx 派生一个带超时/取消的 ctx，设置请求超时时间
	// ctx, cancel := context.WithTimeout(c.Request.Context(), 200*time.Millisecond)
	// defer cancel()
	if err := u.userService.SignUp(ctx, req); err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(err))
		return
	}

	// 3.返回响应
	c.JSON(http.StatusOK, util.JsonRsp("ok"))
}

// LoginHandler
// @Summary 用户登录
// @Tags User
// @Accept json
// @Produce json
// @Param data body dto.LoginReq true "登录请求信息"
// @Success 200 {object} util.JsonResult
// @Router /v1/users/login [post]
func (u *UserHandler) LoginHandler(c *gin.Context) {
	// 获取参数和参数校验
	var req = new(dto.LoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("LoginHandler", zap.Any("LoginReq", req), zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))
		return
	}

	// 登录逻辑处理
	ctx := c.Request.Context()
	res, err := u.userService.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(err))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, util.JsonRsp(res))
}

// LogoutHandler
// @Summary 用户注销
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/users/logout [post]
func (u *UserHandler) LogoutHandler(c *gin.Context) {
	// 获取参数和参数校验
	var req = new(dto.LogoutReq)

	userId, exists := c.Get("UserId")
	if !exists {
		zap.L().Error("UserId not exist")
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.UNAUTHORIZED,
			Msg:  e.GetMsg(e.UNAUTHORIZED),
		}))
	}

	req.UserId = userId.(int64)

	// 注销逻辑处理
	ctx := c.Request.Context()
	err := u.userService.Logout(ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(err))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, util.JsonRsp("ok"))
}
