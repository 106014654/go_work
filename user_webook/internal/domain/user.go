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
	Ctime        time.Time
}
