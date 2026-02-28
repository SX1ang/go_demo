package app

import (
	"demoProject/internal/api"

	_ "demoProject/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// 路由注册
type Router struct {
	basic *api.BasicHandler
	user  *api.UserHandler
}

func NewRouter(basic *api.BasicHandler, user *api.UserHandler) *Router {
	return &Router{
		basic: basic,
		user:  user,
	}
}

func (r *Router) With(engine *gin.Engine) {
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// default
	defaultGroup := engine.Group("")
	defaultGroup.GET("/health", r.basic.Health)

	entityGroup := engine.Group("v1")
	entityGroup.POST("/signup", r.user.SignUpHandler)
}
