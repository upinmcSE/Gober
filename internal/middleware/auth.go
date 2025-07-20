package middleware

import (
	"Gober/pkg/cache"
	"Gober/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	jwtService   jwt.TokenService
	cacheService cache.RedisCacheService
)

func InitAuthMiddleware(token jwt.TokenService, cache cache.RedisCacheService) {
	jwtService = token
	cacheService = cache
}

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

		_, claims, err := jwtService.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})

			return
		}

		if jti, ok := claims["jti"].(string); ok {
			key := "blacklist:" + jti
			exists, err := cacheService.Exists(key)
			if err == nil && exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token revoked",
				})

				return
			}
		}

		payload, err := jwtService.DecryptAccessTokenPayload(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})

			return
		}

		ctx.Set("accountID", payload.AccountID)
		ctx.Set("email", payload.Email)
		ctx.Set("role", payload.Role)

		ctx.Next()
	}
}
