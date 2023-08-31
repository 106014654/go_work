package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go_work/user_webook/internal/web"
)

type LoginJWTMiddleware struct {
	path []string
}

func NewLoginJWTMiddleware() *LoginJWTMiddleware {
	return &LoginJWTMiddleware{}
}

func (lm *LoginJWTMiddleware) AddIngorePath(path string) *LoginJWTMiddleware {
	lm.path = append(lm.path, path)
	return lm
}

func (lm *LoginJWTMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, v := range lm.path {
			if ctx.Request.URL.Path == v {
				return
			}
		}

		//tokenHeader := ctx.GetHeader("Authorization")
		//if tokenHeader == "" {
		//	// 没登录
		//	fmt.Println("tokenHeader")
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//
		//segs := strings.Split(tokenHeader, " ")
		//if len(segs) != 2 {
		//	// 没登录，有人瞎搞
		//	fmt.Println("tokenHeader split")
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//tokenStr := segs[1]

		tokenStr := ctx.GetHeader("x-jwt-token")

		claims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			// 没登录
			fmt.Println("token err1")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			fmt.Println("token nil")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			fmt.Println("token err2")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", claims)
	}
}
