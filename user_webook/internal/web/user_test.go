package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go_work/user_webook/internal/service"
	svcmocks "go_work/user_webook/internal/service/mocks"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserloginJWT(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.CodeServiceInter
		reqBody  string
		wantCode bool
		wantBody Result
	}{
		{
			name: "验证码校验通过",
			mock: func(ctrl *gomock.Controller) service.CodeServiceInter {
				usvc := svcmocks.NewMockCodeServiceInter(ctrl)
				usvc.EXPECT().Verify(gomock.Any(),
					"login", "17312345678", "123456").Return(true, nil)
				return usvc
			},
			reqBody: `{
"biz":"login",
"phone":"17312345678",
"code":"123456"
}`,
			wantCode: true,
			wantBody: Result{
				Code: 0,
				Msg:  "验证码校验通过",
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

			req, err := http.NewRequest(http.MethodPost, "/user/logincode", bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}
