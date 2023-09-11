package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go_work/user_webook/internal/repository"
	"go_work/user_webook/internal/service"
	svcmocks "go_work/user_webook/internal/service/mocks"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSMSCode(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6380",
	})

	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.CodeServiceInter
		reqBody string

		before func(t *testing.T)
		after  func(t *testing.T)

		wantBody Result
	}{
		{
			name: "验证码校验通过",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:123456476", "391051",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:123456476").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "391051", val)
			},
			mock: func(ctrl *gomock.Controller) service.CodeServiceInter {
				usvc := svcmocks.NewMockCodeServiceInter(ctrl)
				usvc.EXPECT().Verify(gomock.Any(),
					"login", "123456476", "391051").Return(true, nil)
				return usvc
			},
			reqBody: `{
"biz":"login",
"phone":"123456476",
"code":"391051"
}`,
			wantBody: Result{
				Code: 0,
				Msg:  "验证码校验通过",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:123456476", "391051",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				_, err := rdb.Del(ctx, "phone_code:login:123456476").Result()
				cancel()
				assert.NoError(t, err)

			},
			mock: func(ctrl *gomock.Controller) service.CodeServiceInter {
				usvc := svcmocks.NewMockCodeServiceInter(ctrl)
				usvc.EXPECT().Verify(gomock.Any(),
					"login", "123456476", "391053").Return(false, repository.ErrCodeVerifyErrUnknownForCode)
				return usvc
			},
			reqBody: `{
"biz":"login",
"phone":"123456476",
"code":"391053"
}`,

			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证码有误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
				_, err := rdb.Set(ctx, "phone_code:login:123456476", "391051",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				_, err := rdb.GetDel(ctx, "phone_code:login:123456476").Result()
				cancel()
				assert.NoError(t, err)

			},
			mock: func(ctrl *gomock.Controller) service.CodeServiceInter {
				usvc := svcmocks.NewMockCodeServiceInter(ctrl)
				usvc.EXPECT().Verify(gomock.Any(),
					"login", "123456476", "391053").Return(false, nil)
				return usvc
			},
			reqBody: `{
"biz":"login",
"phone":"123456476",
"code":"391053"
}`,

			wantBody: Result{
				Code: 4,
				Msg:  "验证码有误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			h := NewUserHandler(nil, tc.mock(ctrl))
			h.RegisteRoute(server)

			req, err := http.NewRequest(http.MethodPost, "/users/logincode", bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
		})
	}
}
