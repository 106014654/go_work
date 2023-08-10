package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleware struct {
	path []string
}

func NewLoginMiddleware() *LoginMiddleware {
	return &LoginMiddleware{}
}

func (lm *LoginMiddleware) AddIngorePath(path string) *LoginMiddleware {
	lm.path = append(lm.path, path)
	return lm
}

func (lm *LoginMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, v := range lm.path {
			if ctx.Request.URL.Path == v {
				return
			}
		}

		sess := sessions.Default(ctx)
		id := sess.Get("user_id")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
