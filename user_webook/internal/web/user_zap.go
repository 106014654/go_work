package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_work/user_webook/internal/domain"
	"net/http"
)

func ZapUserReq[T any](fn func(ctx *gin.Context, req T) (domain.User, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		zap.L().Debug("请求接口开始", zap.String("action", "login"), zap.Any("req", req))
		if err := ctx.Bind(&req); err != nil {
			return
		}
		_, err := fn(ctx, req)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "系统错误")
			zap.L().Debug("请求接口失败", zap.String("action", "login"), zap.Error(err))
		}
		ctx.String(http.StatusOK, "登录成功")
	}
}
