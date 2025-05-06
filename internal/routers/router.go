package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
	r := gin.Default()
	// r.GET("/ping", Pong)

	v1 := r.Group("/v1") 
	{
		v1.GET("/ping", Pong)
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