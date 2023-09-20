package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_work/user_webook/internal/web"
	"go_work/user_webook/internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisteRoute(server)
	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	corsHandler()
	return []gin.HandlerFunc{
		middleware.NewLoginMiddleware().
			AddIngorePath("/users/login").
			AddIngorePath("/users/profile").
			AddIngorePath("/users/logincode").
			AddIngorePath("/users/signup").Build(),
	}
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	})
}
