package dao

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/go-sql-driver/mysql"
)

type UserDAOInter interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
	UpdateSmsCntByPhone(ctx context.Context, phone string, cnt int64) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) UserDAOInter {
	return &GORMUserDAO{
		db: db,
	}
}

func (ud *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error

	return u, err
}

func (ud *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error

	return u, err
}

func (ud *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("id = ?", id).First(&u).Error

	return u, err
}

func (ud *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now

	err := ud.db.WithContext(ctx).Create(&u).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062

		if mysqlErr.Number == uniqueConflictsErrNo {
			return errors.New("邮箱冲突")
		}
	}

	return err
}

func (ud *GORMUserDAO) Update(ctx context.Context, u User) error {
	err := ud.db.WithContext(ctx).Where("id = ?", u.Id).
		Updates(User{NickName: u.NickName, Birth: u.Birth, Introduction: u.Introduction}).Error
	return err
}

func (ud *GORMUserDAO) UpdateSmsCntByPhone(ctx context.Context, phone string, cnt int64) error {
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).
		UpdateColumn("sms_cnt", cnt).Error
	return err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string

	// 往这面加
	NickName     string
	Birth        string
	Introduction string

	SmsCnt int64

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
