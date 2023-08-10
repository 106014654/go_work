package repository

import (
	"context"
	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.GORMUserDAO
}

func NewUserRepository(d *dao.GORMUserDAO) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
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
