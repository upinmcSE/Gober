package initialize

import (
	"Gober/internal/routers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine{
	r := gin.Default()
	// r.GET("/ping", Pong)

	routers.UserRouter(r, db)
	routers.AuthRouter(r, db)

	return r
}
