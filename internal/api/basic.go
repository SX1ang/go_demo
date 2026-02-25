package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasicHandler struct {
}

func NewBasicHandler() BasicHandler {
	return BasicHandler{}
}

// Health
// @Summary 探活接口
// @Produce json
// @Success 200 {string} json "{"msg":"success"}"
// @Router 	/health [get]
func (b *BasicHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
