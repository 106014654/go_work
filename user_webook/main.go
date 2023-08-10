package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"go_work/user_webook/init_web"
	"go_work/user_webook/internal/repository"
	dao2 "go_work/user_webook/internal/repository/dao"
	"go_work/user_webook/internal/service"
	"go_work/user_webook/internal/web"
	"go_work/user_webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//web server setting
	server := init_web.InitWebServer()
	//middleware
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginMiddleware().AddIngorePath("/users/login").
		AddIngorePath("/users/signup").Build())
	//initdb
	db := initDb()
	InitTable(db)
	//init user handler
	uHandle := initHandler(db)
	//register route
	uHandle.RegisteRoute(server)

	server.Run(":9090")
}

func initHandler(db *gorm.DB) *web.UserHandler {
	dao := dao2.NewGORMUserDAO(db)
	rps := repository.NewUserRepository(dao)
	svc := service.NewUserService(rps)
	uHandle := web.NewUserHandler(svc)
	return uHandle
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	return db
}

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&dao2.User{})
}
