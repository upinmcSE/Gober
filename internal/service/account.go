package service

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/repo/mysql"
	"Gober/pkg/cache"
	"Gober/utils/jwt"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateAccountParams struct {
	Email    string
	Password string
}

type CreateAccountOutput struct {
	ID uint64
}

type CreateSessionParams struct {
	Email    string
	Password string
}

type CreateSessionOutput struct {
	Account      *gober.Account
	AccessToken  string
	RefreshToken string
}

type AccountService interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
	RefreshSession(ctx context.Context, refreshToken string) (CreateSessionOutput, error)
	GetAccountByID(ctx context.Context, id uint64) (*gober.Account, error)
	DeleteSession(ctx context.Context, accessToken string, refreshToken string) error
}

type accountService struct {
	db    mysql.AccountDatabase
	token jwt.TokenService
	hash  Hash
	cache cache.RedisCacheService
}

func databaseAccountToProtoAccount(account mysql.Account) *gober.Account {
	return &gober.Account{
		AccountId: account.ID,
		Email:     account.Email,
		Role:      string(account.Role),
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateAccount implements AccountService.
func (a *accountService) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	accountExist, err := a.db.GetAccountByEmail(ctx, params.Email)
	if err != nil {
		return CreateAccountOutput{}, status.Error(codes.Internal, "failed to check account exits")
	}

	if accountExist != nil {
		return CreateAccountOutput{}, status.Error(codes.AlreadyExists, "email is already used")
	}

	hashedPassword, err := a.hash.Hash(ctx, params.Password)
	if err != nil {
		return CreateAccountOutput{}, status.Error(codes.Internal, "failed to hash password")
	}

	account, err := a.db.CreateAccount(ctx, &mysql.Account{
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return CreateAccountOutput{}, status.Error(codes.Internal, "failed to create account")
	}

	return CreateAccountOutput{
		ID: account.ID,
	}, nil

}

// CreateSession implements AccountService.
func (a *accountService) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error) {
	account, err := a.db.GetAccountByEmail(ctx, params.Email)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to get account by email")
	}

	checked, err := a.hash.IsHashEqual(ctx, params.Password, account.Password)

	fmt.Println("account:", account)

	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to check password")
	}
	if !checked {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "password is not correct")
	}

	accessToken, err := a.token.GenerateAccessToken(account)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "failed to generate token")
	}

	refreshToken, err := a.token.GenerateRefreshToken(account)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to generate refresh token")
	}

	if err := a.token.StoreRefreshToken(refreshToken); err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to store refresh token")
	}

	fmt.Println("ouput:", accessToken, refreshToken)

	return CreateSessionOutput{
		Account:      databaseAccountToProtoAccount(*account),
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (a *accountService) RefreshSession(ctx context.Context, refreshTokenStr string) (CreateSessionOutput, error) {
	tokenStr, err := a.token.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
	}

	account, err := a.db.GetAccountByID(ctx, tokenStr.AccountID)
	if err != nil || account == nil {
		return CreateSessionOutput{}, status.Error(codes.NotFound, "failed to get account by id")
	}

	accessToken, err := a.token.GenerateAccessToken(account)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := a.token.GenerateRefreshToken(account)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to generate refresh token")
	}

	if err := a.token.RevokeToken(refreshTokenStr); err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to revoke old refresh token")
	}

	if err := a.token.StoreRefreshToken(refreshToken); err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to store new refresh token")
	}

	fmt.Println("ouput:", accessToken, refreshToken)

	return CreateSessionOutput{
		Account:      databaseAccountToProtoAccount(*account),
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil

}

// GetAccountByID implements AccountService.
func (a *accountService) GetAccountByID(ctx context.Context, id uint64) (*gober.Account, error) {
	account, err := a.db.GetAccountByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get account by id")
	}

	if account == nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	return databaseAccountToProtoAccount(*account), nil
}

func InitAdminAccount(db mysql.AccountDatabase, hash Hash) error {
	// Check if the admin account already exists
	adminAccount, err := db.GetAccountByEmail(context.Background(), "admin@gmail.com")
	if err != nil {
		return status.Error(codes.Internal, "failed to check admin account existence")
	}
	if adminAccount != nil {
		return nil // Admin account already exists, no need to create
	}
	hashPassword, err := hash.Hash(context.Background(), "admin")
	if err != nil {
		return status.Error(codes.Internal, "failed to hash admin password")
	}
	// Create the admin account with a default password
	_, err = db.CreateAccount(context.Background(), &mysql.Account{
		Email:    "admin@gmail.com",
		Password: hashPassword,
		Role:     mysql.Manager,
	})
	if err != nil {
		return status.Error(codes.Internal, "failed to create admin account")
	}
	return nil
}

func (a *accountService) DeleteSession(ctx context.Context, accessToken string, refreshToken string) error {
	_, claims, err := a.token.ParseToken(accessToken)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid access token")
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		a.cache.Set(key, "revoked", ttl)
	}

	_, err = a.token.ValidateRefreshToken(refreshToken)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid or expired refresh token")
	}

	if err := a.token.RevokeToken(refreshToken); err != nil {
		return status.Error(codes.Internal, "failed to revoke refresh token")
	}

	return nil

}

func NewAccountService(db mysql.AccountDatabase, token jwt.TokenService, hash Hash, cache cache.RedisCacheService) AccountService {
	err := InitAdminAccount(db, hash)
	if err != nil {
		return nil
	}

	return &accountService{
		db:    db,
		token: token,
		hash:  hash,
		cache: cache,
	}
}
