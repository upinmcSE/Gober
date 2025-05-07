package initialize

import (
	"Gober/internal/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine{
	r := gin.Default()
	// r.GET("/ping", Pong)

	//r.Use(middlewares.AuthMiddleware()) // Sử dụng middleware xác thực cho tất cả các route

	routers.UserRouter(r)

	return r
}
