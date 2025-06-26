package initialize

import (
	"Gober/configs"
	"Gober/internal/auth/handler/http"
	di "Gober/internal/initialize/wire"
	"Gober/internal/middlewares"
	"Gober/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, config *configs.Config) *gin.Engine {
	// Khởi tạo Gin Engine
	// var r *gin.Engine
	// if gin.Mode() == gin.ReleaseMode {
	// 	r = gin.New()
	// } else {	
	// 	r = gin.Default()
	// }

	r := gin.Default() // Khởi tạo Gin Engine với logging và recovery middleware

	// middlewares
	// r.Use(middlewares.CORS)
	r.Use(middlewares.ValidatorMiddleware())
	r.Use(middlewares.ApiKeyMiddleware(config.Server.ApiKey)) // Middleware kiểm tra API Key

	// r.Use() // logging
	// r.Use() // limiter global

	r.GET("/ping/100", func(ctx *gin.Context) {

		response.SuccessResponse(ctx, "pong")
	})

	// r.GET("/ping/200", response.Wrap(func(ctx *gin.Context) (res interface{}, err error) {
	// 	return "pong", nil
	// }))

	// === Đăng ký routes theo module
	v1 := r.Group("/v1/2025")

	authHandler := di.InitAuth(db)
	http.RegisterAuthRoutes(v1, authHandler)

	// userHandler := initialize.InitUser(db)

	return r
}