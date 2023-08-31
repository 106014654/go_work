package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/redis/go-redis/v9"
	cache2 "go_work/user_webook/internal/repository/cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go_work/user_webook/init_web"
	"go_work/user_webook/internal/repository"
	dao2 "go_work/user_webook/internal/repository/dao"
	"go_work/user_webook/internal/service"
	"go_work/user_webook/internal/web"
	"go_work/user_webook/internal/web/middleware"
)

func main() {

	//web server setting
	server := init_web.InitWebServer()
	//middleware
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginJWTMiddleware().
		AddIngorePath("/users/login").
		AddIngorePath("/users/profile").
		AddIngorePath("/users/signup").Build())
	//initdb
	db := initDb()
	redis := InitRedis()

	//InitTable(db)
	//init user handler
	uHandle := initHandler(db, redis)
	//register route
	uHandle.RegisteRoute(server)

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "success")
	//})
	server.Run(":8081")
}

func initHandler(db *gorm.DB, redis redis.Cmdable) *web.UserHandler {
	dao := dao2.NewGORMUserDAO(db)
	cache := cache2.NewUserCache(redis)
	rps := repository.NewUserRepository(dao, cache)
	svc := service.NewUserService(rps)
	uHandle := web.NewUserHandler(svc)
	return uHandle
}

func initDb() *gorm.DB {
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:30002)/webook"))
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6380",
	})
	return redisClient
}

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&dao2.User{})
}
