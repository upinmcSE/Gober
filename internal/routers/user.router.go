package routers

import (
	c "Gober/internal/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	v1 := r.Group("/v1/user")
	{
		v1.GET("/ping", Pong)
		v1.GET("", c.NewUserController().GetUserById)
	}
}

func Pong(c *gin.Context) {
	name := c.DefaultQuery("name", "upin")
	uid := c.Query("uid")

	c.JSON(http.StatusOK, gin.H{
		"message": "pong " + name,
		"uid":    uid,
	})
}