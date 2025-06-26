package service

import (
	"context"

	appDto "Gober/internal/auth/application/dto"
)

type AuthService interface {
	Login(ctx context.Context, login appDto.AccountAppDTO) (appDto.AccountAppLoginDTO, error)
	Create(ctx context.Context, account appDto.AccountAppDTO) (uint, error)
}