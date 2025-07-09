package service

import (
	"Gober/configs"
	"Gober/internal/repo/mysql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

type tokenService struct {}

type TokenService interface {
	GenerateToken(claims *mysql.Account) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	ExtractEmail(token string) (string, error)
}

// GenerateToken implements TokenService.
func (t *tokenService) GenerateToken(account *mysql.Account) (string, error) {

	config := configs.GetConfig()
	if config == nil {
		return "", errors.New("configuration not found")
	}

	claimsAT := jwt.MapClaims{
		"id":   account.ID,
		"role": account.Role,
		"exp":  config.Security.Expiration, // 1 hour in seconds
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAT)

    return token.SignedString([]byte(config.Security.SecretKey))
}

// ValidateToken implements TokenService.
func (t *tokenService) ValidateToken(tokenString string) (*jwt.Token, error) {

	config := configs.GetConfig()

	if config == nil {
		return nil, errors.New("configuration not found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.Security.SecretKey), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid or expired token")
    }

    return token, nil
}

// ExtractEmail implements TokenService.
func (t *tokenService) ExtractEmail(tokenString string) (string, error) {
	token, err := t.ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return "", errors.New("cannot extract claims")
    }

    return claims.Email, nil
}

func NewTokenService(secretKey string) TokenService {
	return &tokenService{}
}
