package repository

import (
	"Gober/internal/auth/domain/model"
	"context"
)



type AuthRepository interface {
	CreateUser(ctx context.Context, account *model.Account) (*model.Account, error)
}