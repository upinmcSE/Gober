package middlewares

import (
	"Gober/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header
		token := c.GetHeader("Authorization")
		if token == "" {
			response.ErrorResponse(c, response.ErrCodeUnauthorized)
			return
		}

		// Kiểm tra tính hợp lệ của token (ví dụ: JWT)
		if !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isValidToken(token string) bool {
	panic("unimplemented")
}