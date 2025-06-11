package routers

import (
	"Gober/internal/database"
	"Gober/internal/handler"
	"Gober/internal/logic"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRouter(r *gin.Engine, db *gorm.DB) {

	// inject dependencies
	authRepository := database.NewUserDataAccessor(db)
	authService := logic.NewAuthLogic(authRepository, logic.NewHash())
	authController := handler.NewAuthHandler(authService)

	v1 := r.Group("/v1/auth")
	{
		v1.POST("/login", authController.Login)
		v1.POST("/register", authController.Register)
	}
}