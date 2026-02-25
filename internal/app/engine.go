package app

import (
	"demoProject/pkg/util"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(util.GinLogger(), util.GinRecovery(true))
	return engine
}
