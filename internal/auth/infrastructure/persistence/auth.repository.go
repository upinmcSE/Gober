package persistence

import (
	"Gober/internal/auth/domain/model"
	"Gober/internal/auth/domain/repository"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

// CreateUser implements repository.AuthRepository.
func (a *authRepository) CreateUser(ctx context.Context, account *model.Account) (*model.Account, error) {
	createdAccount, err := a.Create(ctx, account)

	if err != nil {
		return nil, err
	}
	if createdAccount == nil {
		return nil, fmt.Errorf("account creation did not return an account instance")
	}

	return createdAccount, nil

}

func (ar *authRepository) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	result := ar.db.WithContext(ctx).Create(account)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create account: %w", result.Error)
	}
	return account, nil
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}
