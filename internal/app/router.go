package app

import (
	"demoProject/internal/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// 路由注册
type Router struct {
	basic api.BasicHandler
}

func NewRouter(basic api.BasicHandler) *Router {
	return &Router{
		basic: basic,
	}
}

func (r *Router) With(engine *gin.Engine) {
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	defaultGroup := engine.Group("")
	defaultGroup.GET("/health", r.basic.Health)

	//entityGroup := engine.Group("v1")
}
