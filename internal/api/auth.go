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

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RenewAccessTokenHandler
// @Summary access token 刷新
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body dto.RenewAccessTokenReq true "刷新access token请求信息"
// @Success 200 {object} util.JsonResult
// @Router /v1/tokens/renew [post]
func (a *AuthHandler) RenewAccessTokenHandler(c *gin.Context) {
	var req = new(dto.RenewAccessTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("RenewAccessTokenHandle", zap.Any("RenewAccessTokenReq", req), zap.Error(err))
		c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
			Code: e.INVALID_PARAMS,
			Msg:  e.GetMsg(e.INVALID_PARAMS),
		}))

		return
	}

	ctx := c.Request.Context()
	res, err := a.authService.RenewAccessToken(ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(err))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, util.JsonRsp(res))

}

// RevokeRefreshTokenHandler
// @Summary refresh token 吊销
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.JsonResult
// @Router /v1/tokens/revoke [post]
func (a *AuthHandler) RevokeRefreshTokenHandler(c *gin.Context) {
	var req = new(dto.RevokeRefreshTokenReq)
	req.UserId = c.MustGet("UserId").(int64)

	ctx := c.Request.Context()
	err := a.authService.RevokeRefreshToken(ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.JsonRsp(err))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, util.JsonRsp("ok"))
}
