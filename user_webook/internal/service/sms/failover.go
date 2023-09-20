package sms

import (
	"context"
	"errors"
	"fmt"
	"go_work/user_webook/internal/service"
	"sync/atomic"
)

/**
与接入服务商提供文档对已知错误码配置成常量 统一管理， 根据错误码处理对应逻辑
优点：对各种情况提前设置好处理方式
缺点：如错误信息有变更，无法第一时间通知，可能会引起短时间内用户无法使用短信服务
*/

var (
	ERRORSMSSERVERSYSTEMFAIL = errors.New("服务器商系统故障")
)

type FailoverSMSService struct {
	svcs []SMSService
	idx  int32

	cnt int32

	maxCnt int32

	usvc service.UserServiceInter
}

func (f *FailoverSMSService) getUserSmsCnt(ctx context.Context, phone string) int64 {
	user, err := f.usvc.GetUserInfoByPhone(ctx, phone)

	if err != nil {
		return 0
	}
	return user.SmsCnt
}

func (f *FailoverSMSService) SendCode(ctx context.Context, biz string, phone string) error {
	idx := atomic.LoadInt32(&f.idx)
	cnt := f.getUserSmsCnt(ctx, phone)
	svc := f.svcs[idx]

	if int32(cnt) > f.maxCnt { //最大尝试次数
		newIdx := (idx + 1) % int32(len(f.svcs))
		if atomic.CompareAndSwapInt32(&f.idx, idx, newIdx) {
			atomic.StoreInt32(&f.cnt, 0)
			err := f.usvc.EditSmsCntByPhone(ctx, phone, 0)
			if err != nil {
				return err
			}
		}
		idx = atomic.LoadInt32(&f.idx)
	}
	if cnt != 0 {
		var er error
		go func() {
			err := svc.SendCode(ctx, biz, phone)
			if err != nil {
				atomic.StoreInt32(&f.cnt, int32(cnt)+1)
				_ = f.usvc.EditSmsCntByPhone(ctx, phone, cnt+1)
			} else {
				atomic.StoreInt32(&f.cnt, 0)
				_ = f.usvc.EditSmsCntByPhone(ctx, phone, 0)
			}
			er = err
		}()
		return er
	}

	err := svc.SendCode(ctx, biz, phone)

	cntAdd := cnt + 1
	switch err {
	case ERRORSMSSERVERSYSTEMFAIL:
		f.cnt += 1
		fmt.Printf("cnt:%d", f.cnt)
		_ = f.usvc.EditSmsCntByPhone(ctx, phone, cntAdd)
		if f.idx != int32(len(f.svcs)) {
			_ = svc.SendCode(ctx, biz, phone)
		} else {
			return err
		}
	case nil:
		atomic.StoreInt32(&f.cnt, 0)
		_ = f.usvc.EditSmsCntByPhone(ctx, phone, 0)
		return nil
	default:
		return err
	}
	return nil
}

func NewFailoverSMSService(svc []SMSService, idx, cnt, mnt int32, rsp service.UserServiceInter) *FailoverSMSService {
	return &FailoverSMSService{
		svcs:   svc,
		idx:    idx,
		cnt:    cnt,
		maxCnt: mnt,
		usvc:   rsp,
	}
}
