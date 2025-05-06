package routers

import (
	c "Gober/internal/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
	r := gin.Default()
	// r.GET("/ping", Pong)

	//r.Use(middlewares.AuthMiddleware()) // Sử dụng middleware xác thực cho tất cả các route

	v1 := r.Group("/v1/user") 
	{
		v1.GET("/ping", Pong)
		v1.GET("", c.NewUserController().GetUserById)
	}

	return r
}

func Pong(c *gin.Context) {
	name := c.DefaultQuery("name", "upin")
	uid := c.Query("uid")
  
	c.JSON(http.StatusOK, gin.H{
	  "message": "pong " + name,
	  "uid": uid,
	})
}