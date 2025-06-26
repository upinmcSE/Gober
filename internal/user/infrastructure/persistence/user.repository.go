package persistence

import (
	"Gober/internal/user/domain/model"
	"Gober/internal/user/domain/repository"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// EmailExists implements repository.UserRepository.
func (u *userRepository) EmailExists(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	result := u.db.WithContext(ctx).Where("email = ?", email).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Email does not exist
		}
		return nil, fmt.Errorf("failed to check email existence: %w", result.Error)
	}

	return &account, nil
}

// GetProfile implements repository.UserRepository.
func (u *userRepository) GetProfile(ctx context.Context, userID uint) (*model.Account, error) {
	var account model.Account
	result := u.db.WithContext(ctx).Where("id = ?", userID).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // User profile not found
		}
		return nil, fmt.Errorf("failed to get user profile: %w", result.Error)
	}

	return &account, nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}
