package app

import (
	"context"
	"demoProject/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine  *gin.Engine
	router  *Router
	address string
}

func NewServer(engine *gin.Engine, router *Router, cfg *config.Config) *Server {
	return &Server{
		engine:  engine,
		router:  router,
		address: fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port),
	}
}

func (s *Server) Run() error {
	// 注册路由
	s.router.With(s.engine)

	// 启动 HTTP 服务
	srv := &http.Server{
		Addr:    s.address,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: %s\n", zap.Error(err))
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down server...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("Server exiting")
	return nil
}
