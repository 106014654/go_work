package service

import (
	"context"
	"errors"
	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidEmailOrPassword = errors.New("邮箱或密码对")

type UserService struct {
	ur *repository.UserRepository
}

func NewUserService(usr *repository.UserRepository) *UserService {
	return &UserService{
		ur: usr,
	}
}

func (user *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)
	return user.ur.Create(ctx, u)
}

func (user *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	us, err := user.ur.FindByEmail(ctx, email)
	if err == gorm.ErrRecordNotFound {
		return domain.User{}, nil
	}

	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidEmailOrPassword
	}

	return us, nil
}

func (user *UserService) EditUserDetail(ctx context.Context, id int64, name, birth, intro string) error {
	err := user.ur.EditByUserId(ctx, domain.User{
		Id:           id,
		NickName:     name,
		Birth:        birth,
		Introduction: intro,
	})

	return err
}

func (user *UserService) GetUserInfo(ctx context.Context, id int64) (domain.User, error) {
	return user.ur.FindByUserId(ctx, id)
}