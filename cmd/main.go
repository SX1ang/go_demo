//go:build !wireinject
// +build !wireinject

package main

import (
	"demoProject/config"
	"demoProject/pkg/util"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	util.InitLogger(cfg)
	defer zap.L().Sync()

	server, err := InitServer(cfg)
	if err != nil {
		panic("init server err")
	}

	server.Run()

}
