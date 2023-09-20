package sms

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/service"
	svcmocks "go_work/user_webook/internal/service/mocks"
	smsmocks "go_work/user_webook/internal/service/sms/mocks"
	"testing"
)

func TestFailSMS(t *testing.T) {
	testCases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) []SMSService
		mock1 func(ctrl *gomock.Controller) service.UserServiceInter
		idx   int32

		cnt int32

		maxCnt int32

		wantIdx int32
		wantCnt int32

		wantError error
	}{
		{
			name: "成功发送",
			mock: func(ctrl *gomock.Controller) []SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc1 := smsmocks.NewMockSMSService(ctrl)

				svc0.EXPECT().SendCode(gomock.Any(), "login", "123445678").
					Return(nil)
				return []SMSService{svc0, svc1}
			},
			mock1: func(ctrl *gomock.Controller) service.UserServiceInter {
				usvc := svcmocks.NewMockUserServiceInter(ctrl)
				usvc.EXPECT().GetUserInfoByPhone(gomock.Any(), "123445678").
					Return(domain.User{
						SmsCnt: 0,
					}, nil)
				usvc.EXPECT().EditSmsCntByPhone(gomock.Any(), "123445678", int64(0)).Return(nil)
				return usvc
			},
			idx:     0,
			cnt:     0,
			maxCnt:  3,
			wantIdx: 0,
			wantCnt: 0,
		},
		{
			name: "第一次发送失败，第二次成功",
			mock: func(ctrl *gomock.Controller) []SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc1 := smsmocks.NewMockSMSService(ctrl)

				svc0.EXPECT().SendCode(gomock.Any(), "login", "123445678").
					Return(ERRORSMSSERVERSYSTEMFAIL).AnyTimes()

				return []SMSService{svc0, svc1}
			},
			mock1: func(ctrl *gomock.Controller) service.UserServiceInter {
				usvc := svcmocks.NewMockUserServiceInter(ctrl)

				usvc.EXPECT().GetUserInfoByPhone(gomock.Any(), "123445678").
					Return(domain.User{
						SmsCnt: 0,
					}, nil)

				usvc.EXPECT().EditSmsCntByPhone(gomock.Any(), "123445678", int64(1)).Return(nil)

				return usvc
			},
			idx:       0,
			cnt:       0,
			maxCnt:    3,
			wantIdx:   0,
			wantCnt:   1,
			wantError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailoverSMSService(tc.mock(ctrl), tc.idx, tc.cnt, tc.maxCnt, tc.mock1(ctrl))

			err := svc.SendCode(context.Background(), "login", "123445678")

			assert.Equal(t, tc.wantIdx, svc.idx)
			assert.Equal(t, tc.wantCnt, svc.cnt)
			assert.Equal(t, tc.wantError, err)

		})
	}
}
