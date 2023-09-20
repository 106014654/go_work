package main

import (
	"go_work/user_webook/config"
	"go_work/user_webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//InitTable()
	server := InitWebServer()
	server.Run(":8081")
}

func InitTable() error {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	return db.AutoMigrate(&dao.User{})
}
