package repository

import (
	"Gober/internal/user/domain/model"
	"context"
)

type UserRepository interface {
	EmailExists(ctx context.Context, email string) (*model.Account, error)
	GetProfile(ctx context.Context, userID uint) (*model.Account, error)
}