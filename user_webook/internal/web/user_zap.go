package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_work/user_webook/internal/domain"
	"net/http"
)

func ZapUserReq(us *UserHandler, ctx *gin.Context, email, password string) (domain.User, error) {
	zap.L().Debug("请求接口开始", zap.String("action", "login"), zap.String("param", email))
	user, err := us.uservice.Login(ctx, email, password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		zap.L().Debug("请求接口失败", zap.String("action", "login"), zap.Error(err))
	}
	return user, err
}
