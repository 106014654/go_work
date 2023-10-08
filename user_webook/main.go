package main

import (
	"fmt"
	"go_work/user_webook/config"
	"go_work/user_webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
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

func quick_sort(nums []int, l, r int) {
	if l >= r {
		return
	}
	rand.Seed(time.Now().Unix())
	p := rand.Intn(r-l+1) + l
	nums[r], nums[p] = nums[p], nums[r]
	i := l - 1
	for j := l; j < r; j++ {
		if nums[j] < nums[r] {
			j++
			nums[i], nums[j] = nums[j], nums[i]
			fmt.Println(nums, i, j)
		}
	}

	nums[i+1], nums[r] = nums[r], nums[i+1]
	quick_sort(nums, l, i-1)
	quick_sort(nums, i+1, r)
}
