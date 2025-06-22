package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

type tokenService struct {
    secretKey string
}

type TokenService interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	ExtractEmail(token string) (string, error)
}

// GenerateToken implements TokenService.
func (t *tokenService) GenerateToken(claims jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(t.secretKey))
}

// ValidateToken implements TokenService.
func (t *tokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(t.secretKey), nil
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
	return &tokenService{
        secretKey: secretKey,
    }
}
