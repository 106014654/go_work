//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go_work/user_webook/internal/repository"
	"go_work/user_webook/internal/repository/cache"
	"go_work/user_webook/internal/repository/dao"
	"go_work/user_webook/internal/service"
	"go_work/user_webook/internal/web"
	"go_work/user_webook/ioc"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		ioc.InitRedis, ioc.InitDb,

		dao.NewGORMUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		service.NewCodeService,
		service.NewUserService,

		web.NewUserHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return new(gin.Engine)
}
