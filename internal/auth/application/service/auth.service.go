package service

import (
	"context"

	appDto "Gober/internal/auth/application/dto"
	hlerDto "Gober/internal/auth/handler/dto"
)

type AuthService interface {
	Login(ctx context.Context, login hlerDto.AccountLoginReq) (hlerDto.AccountLoginRes, error)
	Create(ctx context.Context, account appDto.AccountAppDTO) (uint, error)
}