//go:build wireinject
// +build wireinject

package main

import (
	"demoProject/config"
	"demoProject/internal/api"
	"demoProject/internal/app"
	"demoProject/internal/infra"
	"demoProject/internal/infra/mysql"
	"demoProject/internal/infra/redis"
	"demoProject/internal/repo"
	"demoProject/internal/service"

	"github.com/google/wire"
)

func InitConfig() (*config.Config, error) {
	wire.Build(
		config.InitConfig,
	)

	return nil, nil
}

func InitServer(cfg *config.Config) (*app.Server, func(), error) {
	wire.Build(

		// db
		mysql.InitDB,

		// redis
		redis.InitRedis,

		// resources
		infra.NewResources,

		// gin-engine & router
		app.NewEngine,
		app.NewRouter,

		// handler
		api.NewBasicHandler,
		api.NewUserHandler,

		// service
		service.NewUserService,

		// repo
		repo.NewMysqlUserRespository,

		// server
		app.NewServer,
	)

	return nil, nil, nil
}
