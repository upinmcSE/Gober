package wire

import (
	"Gober/configs"
	"Gober/internal/auth/application/service"
	"Gober/internal/auth/handler/http"
	authRepo "Gober/internal/auth/infrastructure/persistence"
	userRepo "Gober/internal/user/infrastructure/persistence"

	"gorm.io/gorm"
)

// initializes service, repository, handler, hash, token for auth
func InitAuth(db *gorm.DB) *http.AuthHandler {
	cfg := configs.GetConfig()
	hash := service.NewHashService()
	token := service.NewTokenService(cfg.Security.SecretKey)

	authRepo := authRepo.NewAuthRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	service := service.NewAuthService(authRepo, userRepo, hash, token)
	handler := http.NewAuthHandler(service)
	return handler
}