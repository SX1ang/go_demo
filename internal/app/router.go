package app

import (
	"demoProject/internal/api"
	"demoProject/internal/middleware"
	"demoProject/pkg/util"

	_ "demoProject/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// 路由注册
type Router struct {
	basic      *api.BasicHandler
	user       *api.UserHandler
	auth       *api.AuthHandler
	tokenMaker *util.JWTMaker
}

func NewRouter(
	basic *api.BasicHandler,
	user *api.UserHandler,
	auth *api.AuthHandler,
	tokenMaker *util.JWTMaker) *Router {
	return &Router{
		basic:      basic,
		user:       user,
		auth:       auth,
		tokenMaker: tokenMaker,
	}
}

func (r *Router) With(engine *gin.Engine) {
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// default
	defaultGroup := engine.Group("/")
	defaultGroup.GET("/health", r.basic.Health)

	versionGroup := engine.Group("/v1")

	userGroup := versionGroup.Group("/user")
	{
		userGroup.POST("/signup", r.user.SignUpHandler)
		userGroup.POST("/login", r.user.LoginHandler)

		// 需要登录认证后才能访问 logout
		authUserGroup := userGroup.Group("")
		authUserGroup.Use(middleware.JWT(r.tokenMaker))
		{
			authUserGroup.POST("/logout", r.user.LogoutHandler)
		}
	}

	tokenGroup := versionGroup.Group("/token")
	{
		//  renew 不需要认证 access token
		tokenGroup.POST("/renew", r.auth.RenewAccessTokenHandle)

		// revoke 需要 access token 鉴权
		authTokenGroup := tokenGroup.Group("")
		authTokenGroup.Use(middleware.JWT(r.tokenMaker))
		{
			authTokenGroup.POST("/revoke", r.auth.RevokeRefreshTokenHandle)
		}
	}
}
