package service

import (
	"Gober/configs"
	"Gober/internal/repo/mysql"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type tokenService struct{}

type TokenService interface {
	GenerateToken(claims *mysql.Account) (string, error)
	ValidateToken(token string) (bool, error)
	ExtractEmail(token string) (string, error)
	ExtractAccountID(token string) (string, error)
}

func (t *tokenService) ExtractAccountID(tokenString string) (string, error) {
	config := configs.GetConfig()
	if config == nil {
		return "", status.Error(codes.Internal, "Configuration not found")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Security.SecretKey), nil
	})

	if err != nil {
		return "", status.Error(codes.Unauthenticated, "Token is invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", status.Error(codes.Internal, "Lỗi server")
	}

	return fmt.Sprintf("%v", claims["id"]), nil
}

// GenerateToken implements TokenService.
func (t *tokenService) GenerateToken(account *mysql.Account) (string, error) {

	config := configs.GetConfig()
	if config == nil {
		return "", status.Error(codes.Internal, "Configuration not found")
	}

	claimsAT := jwt.MapClaims{
		"id":   account.ID,
		"role": account.Role,
		"exp":  time.Now().Add(time.Duration(config.Security.Expiration) * time.Second).Unix(), // 1 hour in seconds
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAT)

	return token.SignedString([]byte(config.Security.SecretKey))
}

// ValidateToken implements TokenService.
func (t *tokenService) ValidateToken(tokenString string) (bool, error) {
	config := configs.GetConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Security.SecretKey), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token claims:", claims)
		return true, nil
	}

	return false, errors.New("invalid token")
}

// ExtractEmail implements TokenService.
func (t *tokenService) ExtractEmail(tokenString string) (string, error) {
	config := configs.GetConfig()
	if config == nil {
		return "", status.Error(codes.Internal, "Lỗi server")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Security.SecretKey), nil
	})
	if !token.Valid || err != nil {
		return "", status.Error(codes.Unauthenticated, "Token is invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", status.Error(codes.Internal, "Lỗi server")
	}

	return claims.Email, nil
}

func NewTokenService() TokenService {
	return &tokenService{}
}
