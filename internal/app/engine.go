package app

import (
	"demoProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinLogger())
	engine.Use(middleware.GinRecovery(true))
	return engine
}
