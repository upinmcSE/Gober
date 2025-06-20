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

// EmailExists implements repository.AuthRepository.
func (a *authRepository) EmailExists(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	result := a.db.WithContext(ctx).Where("email = ?", email).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Email does not exist
		}
		return nil, fmt.Errorf("failed to check email existence: %w", result.Error)
	}

	return &account, nil
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}
