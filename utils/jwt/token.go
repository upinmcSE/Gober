package jwt

import (
	"Gober/configs"
	"Gober/internal/repo/mysql"
	"Gober/pkg/cache"
	"Gober/utils/crypt"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type EncryptedPayload struct {
	AccountID uint64 `json:"account_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	AccountID uint64    `json:"account_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

type tokenService struct {
	cache  cache.RedisCacheService
	config *configs.Config
}

func (t tokenService) GenerateAccessToken(account *mysql.Account) (string, error) {
	payload := &EncryptedPayload{
		AccountID: account.ID,
		Email:     account.Email,
		Role:      string(account.Role),
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := crypt.EncryptAES(rawData, crypt.JwtEncryptionKey)
	if err != nil {
		return "", err
	}

	ttl := time.Duration(t.config.Security.Expiration.RefreshToken) * time.Second

	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "gober",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.Security.SecretKey))
}

func (t tokenService) GenerateRefreshToken(account *mysql.Account) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		AccountID: account.ID,
		ExpiresAt: time.Now().Add(time.Duration(t.config.Security.Expiration.RefreshToken) * time.Second),
		Revoked:   false,
	}, nil
}

func (t tokenService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	jwtSecret := []byte(t.config.Security.SecretKey)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, status.Error(codes.Unauthenticated, "Failed to parse token claims")
	}

	return token, claims, nil
}

func (t tokenService) DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := t.ParseToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Invalid token data")
	}

	decryptedBytes, err := crypt.DecryptAES(encryptedData, crypt.JwtEncryptionKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(decryptedBytes, &payload); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to parse token payload")
	}

	return &payload, nil
}

func (t tokenService) StoreRefreshToken(token RefreshToken) error {
	ttl := time.Duration(t.config.Security.Expiration.RefreshToken) * time.Second
	cacheKey := "refresh_token:" + token.Token
	return t.cache.Set(cacheKey, token, ttl)
}

func (t tokenService) ValidateRefreshToken(tokenString string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + tokenString

	var refreshToken RefreshToken
	err := t.cache.Get(cacheKey, &refreshToken)

	if err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, status.Error(codes.Unauthenticated, "Invalid or expired refresh token")
	}

	return refreshToken, nil
}

func (t tokenService) RevokeToken(token string) error {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := t.cache.Get(cacheKey, &refreshToken)
	if err != nil {
		return status.Error(codes.NotFound, "Refresh token not found")
	}

	refreshToken.Revoked = true

	return t.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}

func NewTokenService(cache cache.RedisCacheService, config *configs.Config) TokenService {
	return &tokenService{
		cache:  cache,
		config: config,
	}
}
