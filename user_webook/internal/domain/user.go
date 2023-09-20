package domain

import (
	"time"
)

type User struct {
	Id           int64
	Email        string
	Phone        string
	Password     string
	NickName     string
	Birth        string
	Introduction string
	SmsCnt       int64
	Ctime        time.Time
}
