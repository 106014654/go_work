package service

import (
	"context"
	"fmt"
	"go_work/user_webook/internal/repository"
	"math/rand"
)

var (
	ErrCodeSendTooMany = repository.ErrCodeSendTooMany
)

type CodeServiceInter interface {
	SendCode(ctx context.Context,
		biz string, phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}

type CodeService struct {
	repo repository.CodeRepositoryInter
}

func NewCodeService(repo repository.CodeRepositoryInter) CodeServiceInter {
	return &CodeService{
		repo: repo,
	}
}

func (cs *CodeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001

	return fmt.Sprintf("%06d", num)
}

func (cs *CodeService) SendCode(ctx context.Context, biz, phone string) error {
	code := cs.generateCode()
	fmt.Println(code)
	err := cs.repo.Store(ctx, biz, phone, code)

	return err
}

func (cs *CodeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return cs.repo.Verify(ctx, biz, phone, code)
}
