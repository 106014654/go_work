package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/repository/cache"
	"go_work/user_webook/internal/repository/dao"
)

type UserRepository struct {
	dao   *dao.GORMUserDAO
	cache *cache.RedisUserCache
}

func NewUserRepository(d *dao.GORMUserDAO, cache *cache.RedisUserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: cache,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), nil
}

func (ur *UserRepository) FindByUserId(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birth:        u.Birth,
		Introduction: u.Introduction,
	}, nil
}

func (ur *UserRepository) EditByUserId(ctx context.Context, u domain.User) error {
	return ur.dao.Update(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birth:        u.Birth,
		Introduction: u.Introduction,
	})
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//cache
	us, err := ur.cache.Get(ctx, id)

	if err == cache.ERRREDISNIL {
		//dao
		udao, er := ur.dao.FindById(ctx, id)
		fmt.Println("find by mysql")
		if er != nil {
			return domain.User{}, nil
		}
		//set cache
		us = ur.entityToDomain(udao)
		err = ur.cache.Set(ctx, us)
		if err != nil {
			log.Fatalln(err)
		}
		return us, err
	}
	return us, err
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
