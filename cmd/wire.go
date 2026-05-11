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
	"demoProject/pkg/util"

	"github.com/google/wire"
)

func InitConfig() (*config.Config, error) {
	wire.Build(
		config.InitConfig,
	)

	return nil, nil
}

func InitServer(*config.Config) (*app.Server, func(), error) {
	wire.Build(
		// env config
		config.InitEnvConfig,
		config.ProvideSecretKey,

		// jwt
		util.NewJWTMaker,

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
		api.NewAuthHandler,
		api.NewcommunityHandler,

		// service
		service.NewUserService,
		service.NewAuthService,
		service.NewCommunityService,

		// repo
		repo.NewMysqlUserRespository,
		repo.NewMysqlSessionRespository,
		repo.NewMysqlCommunityRespository,

		// server
		app.NewServer,
	)

	return nil, nil, nil
}
