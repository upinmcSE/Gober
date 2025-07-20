package jwt

import (
	"Gober/internal/repo/mysql"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(account *mysql.Account) (string, error)
	GenerateRefreshToken(account *mysql.Account) (RefreshToken, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	StoreRefreshToken(token RefreshToken) error
	ValidateRefreshToken(tokenString string) (RefreshToken, error)
	RevokeToken(token string) error
}
