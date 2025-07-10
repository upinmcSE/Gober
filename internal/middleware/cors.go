package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("CORS Origin:", ctx.Request.Header.Get("Origin"))
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
