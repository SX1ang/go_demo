//go:build !wireinject
// +build !wireinject

package main

import (
	"demoProject/config"
	"demoProject/pkg/snowflake"
	"demoProject/pkg/util"

	"go.uber.org/zap"
)

// @title Go Web Demo Project's API
// @version 1.0
// @description go web demo project

// @contact.name yesongxi
// @contact.email yinshouxiang.email@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

// @tag.name User
// @tag.description 用户相关接口

func main() {

	snowflake.Init("2026-01-01", 1)

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	util.InitLogger(cfg)
	defer zap.L().Sync()

	server, close, err := InitServer(cfg)
	if err != nil {
		panic("init server err")
	}
	defer close()

	server.Run()

}
