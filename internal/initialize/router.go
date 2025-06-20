package initialize

import (
	"Gober/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	// Khởi tạo Gin Engine
	// var r *gin.Engine
	// if gin.Mode() == gin.ReleaseMode {
	// 	r = gin.New()
	// } else {	
	// 	r = gin.Default()
	// }

	r := gin.New()

	// middlewares
	//r.Use(middleware.CORS) // cross
	//r.Use(middleware.ValidatorMiddleware())
	// r.Use() // logging
	// r.Use() // limiter global
	// r.Use(middlewares.Validator())

	r.GET("/ping/100", func(ctx *gin.Context) {

		response.SuccessResponse(ctx, "pong")
	})

	// r.GET("/ping/200", response.Wrap(func(ctx *gin.Context) (res interface{}, err error) {
	// 	return "pong", nil
	// }))

	// === Đăng ký routes theo module
	//v1 := r.Group("/v1/2025")

	//authHandler := initialize.InitAuth(db)
	//http.RegisterAuthRoutes(v1, authHandler)

	// userHandler := initialize.InitUser(db)

	return r
}