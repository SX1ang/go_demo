package app

import (
	"demoProject/config"
	"demoProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewEngine(cfg *config.Config) *gin.Engine {
	if cfg.App.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(middleware.GinLogger())
	engine.Use(middleware.GinRecovery(true))
	return engine
}
