package middleware

import (
	"Gober/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApikeyMiddleware() gin.HandlerFunc {
	config := configs.GetConfig()
	expectedKey := config.Server.ApiKey
	if expectedKey == "" {
		expectedKey = "secret-key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-API-Key"})
			return
		}

		if apiKey != expectedKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		ctx.Next()
	}
}
