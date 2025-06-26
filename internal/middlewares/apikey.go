package middlewares

import (
	"Gober/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware(expectedKey string) gin.HandlerFunc {
	if expectedKey == "" {
		expectedKey = "secret-key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("x-api-key")
		if apiKey == "" {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "API Key is required", "Missing API Key in request header")
			ctx.Abort()
			return
		}

		if apiKey != expectedKey {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid API Key", "The provided API Key is invalid or does not match the expected key")
			ctx.Abort()
			return
		}

		ctx.Set("username", "quoctuan")

		ctx.Next()
	}
}