package routers

import (
	"Gober/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRouter(r *gin.Engine, db *gorm.DB) {

	// tiêm phụ thuộc (DI) thủ công


	v1 := r.Group("/v1/user")
	{
		v1.Use(middlewares.AuthMiddleware())
		v1.PATCH("")
	}
}
