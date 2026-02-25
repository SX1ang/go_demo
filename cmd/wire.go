//go:build wireinject
// +build wireinject

package main

import (
	"demoProject/config"
	"demoProject/internal/api"
	"demoProject/internal/app"

	"github.com/google/wire"
)

func InitConfig() (*config.Config, error) {
	wire.Build(
		config.InitConfig,
	)

	return nil, nil
}

func InitServer(cfg *config.Config) (*app.Server, error) {
	wire.Build(
		
		//// db
		//mysql.InitDB,
		//
		//// redis
		//redis.InitRedis,

		// resources
		//app.NewResources,

		// repo

		// service

		// handler
		api.NewBasicHandler,

		// gin-engine & router
		app.NewEngine,
		app.NewRouter,

		// server
		app.NewServer,
	)

	return nil, nil
}
