package middleware

import (
	"Gober/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		tokenService := service.NewTokenService()
		var valid, err = tokenService.ValidateToken(token)
		if !valid || err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token1"})
			return
		}

		accountID, err := tokenService.ExtractAccountID(token)
		if err != nil && accountID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token2"})
			return
		}

		id, err := strconv.ParseUint(accountID, 10, 64)

		ctx.Set("accountID", id)

		ctx.Next()
	}
}
