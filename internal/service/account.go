package service

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/repo/mysql"
	"context"

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
	Account *gober.Account
	Token   string
}

type AccountService interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
	GetAccountByID(ctx context.Context, id uint64) (*gober.Account, error)
}

type accountService struct {
	db    mysql.AccountDatabase
	token TokenService
	hash  Hash
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
		return CreateAccountOutput{}, err
	}

	return CreateAccountOutput{
		ID: account.ID,
	}, nil

}

// CreateSession implements AccountService.
func (a *accountService) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error) {
	accountExist, err := a.db.GetAccountByEmail(ctx, params.Email)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to check account exits")
	}

	if accountExist == nil {
		return CreateSessionOutput{}, status.Error(codes.NotFound, "email is not found")
	}

	account, err := a.db.GetAccountByEmail(ctx, params.Email)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to get account by email")
	}

	checked, err := a.hash.IsHashEqual(ctx, params.Password, account.Password)

	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to check password")
	}
	if !checked {
		return CreateSessionOutput{}, status.Error(codes.Unauthenticated, "password is not correct")
	}

	token, err := a.token.GenerateToken(account)
	if err != nil {
		return CreateSessionOutput{}, status.Error(codes.Internal, "failed to generate token")
	}

	return CreateSessionOutput{
		Account: databaseAccountToProtoAccount(*account),
		Token:   token,
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

func NewAccountService(db mysql.AccountDatabase, token TokenService, hash Hash) AccountService {
	return &accountService{
		db:    db,
		token: token,
		hash:  hash,
	}
}
