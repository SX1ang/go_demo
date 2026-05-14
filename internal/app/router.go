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
	post      *api.PostHandler
}

func NewRouter(
	cfg *config.Config,
	basic *api.BasicHandler,
	user *api.UserHandler,
	auth *api.AuthHandler,
	tokenMaker *util.JWTMaker,
	community *api.CommunityHandler,
	post *api.PostHandler,
) *Router {
	return &Router{
		cfg:        cfg,
		basic:      basic,
		user:       user,
		auth:       auth,
		tokenMaker: tokenMaker,

		community: community,
		post:      post,
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

	postGroup := versionGroup.Group("/posts", r.jwtMiddleware())
	{
		postGroup.POST("", r.post.CreatePostHandler)
		// 通过http:x.x.x.x/post/4567839465372839487（帖子ID）获取帖子详情
		postGroup.GET("/:id", r.post.PostDetailHandler)

		postGroup.GET("", r.post.PostListHandler)
	}
}
