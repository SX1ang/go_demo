package middleware

import (
	"demoProject/internal/e"
	"demoProject/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWT(tokenMaker *util.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
				Code: e.UNAUTHORIZED,
				Msg:  e.GetMsg(e.UNAUTHORIZED),
			}))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
				Code: e.INVALID_PARAMS,
				Msg:  "Authorization 格式错误，应为 Bearer <token>",
			}))
			c.Abort()
			return
		}

		claims, err := tokenMaker.VerifyToken(parts[1])
		if err != nil || claims.TokenType != util.TokenTypeAccessToken {
			c.JSON(http.StatusOK, util.JsonRsp(&util.CustomizedErr{
				Code: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
				Msg:  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
			}))
			c.Abort()
			return
		}

		// 将当前请求的UserId信息保存到请求的上下文c上
		c.Set("UserId", claims.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("UserId")来获取当前请求的用户信息
	}
}
