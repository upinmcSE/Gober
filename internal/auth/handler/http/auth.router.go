package http

import (
	"Gober/pkg/response"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")

	auth.POST("/register", response.Wrap(handler.RegisterUser))
	auth.POST("/login", response.Wrap(handler.LoginUser))

}