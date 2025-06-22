package di

import (
	"Gober/configs"
	"Gober/internal/auth/application/service"
	"Gober/internal/auth/handler/http"
	"Gober/internal/auth/infrastructure/persistence"

	"gorm.io/gorm"
)

// initializes service, repository, handler, hash, token for auth
func InitAuth(db *gorm.DB) *http.AuthHandler {
	cfg := configs.GetConfig()
	hash := service.NewHashService()
	token := service.NewTokenService(cfg.Security.SecretKey)

	authRepo := persistence.NewAuthRepository(db)
	service := service.NewAuthService(authRepo, hash, token)
	handler := http.NewAuthHandler(service)
	return handler
}