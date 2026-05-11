package app

import (
	"demoProject/config"
	_ "demoProject/docs"
	"demoProject/internal/api"
	"demoProject/internal/middleware"
	"demoProject/pkg/util"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// 路由注册
type Router struct {
	cfg        *config.Config
	basic      *api.BasicHandler
	user       *api.UserHandler
	auth       *api.AuthHandler
	tokenMaker *util.JWTMaker

	community *api.CommunityHandler
}

func NewRouter(
	cfg *config.Config,
	basic *api.BasicHandler,
	user *api.UserHandler,
	auth *api.AuthHandler,
	tokenMaker *util.JWTMaker,
	community *api.CommunityHandler,
) *Router {
	return &Router{
		cfg:        cfg,
		basic:      basic,
		user:       user,
		auth:       auth,
		tokenMaker: tokenMaker,

		community: community,
	}
}

// 封装一个条件中间件
func (r *Router) jwtMiddleware() gin.HandlerFunc {
	if r.cfg.App.Mode == "dev" {
		return func(c *gin.Context) { c.Next() } // dev 环境直接放行
	}
	return middleware.JWT(r.tokenMaker)
}

func (r *Router) With(engine *gin.Engine) {
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// default
	defaultGroup := engine.Group("/")
	defaultGroup.GET("/health", r.basic.Health)

	versionGroup := engine.Group("/v1")

	userGroup := versionGroup.Group("/users")
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

	tokenGroup := versionGroup.Group("/tokens")
	{
		//  renew 不需要认证 access token
		tokenGroup.POST("/renew", r.auth.RenewAccessTokenHandler)

		// revoke 需要 access token 鉴权
		authTokenGroup := tokenGroup.Group("")
		authTokenGroup.Use(middleware.JWT(r.tokenMaker))
		{
			authTokenGroup.POST("/revoke", r.auth.RevokeRefreshTokenHandler)
		}
	}

	communityGroup := versionGroup.Group("/communities", r.jwtMiddleware())
	{
		communityGroup.GET("", r.community.CommunityListHandler)
		communityGroup.GET("/:id", r.community.CommunityDetailHandler)
	}
}
