package grpc

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/service"
	"context"
)

type AccountHandler struct {
	gober.UnimplementedGoberServiceServer
	AccountService service.AccountService
}

func (h *AccountHandler) CreateAccount(ctx context.Context, req *gober.CreateAccountRequest) (*gober.CreateAccountResponse, error) {
	params := service.CreateAccountParams{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.AccountService.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &gober.CreateAccountResponse{
		AccountId: output.ID,
	}, nil
}

func (h *AccountHandler) CreateSession(ctx context.Context, req *gober.CreateSessionRequest) (*gober.CreateSessionResponse, error) {
	params := service.CreateSessionParams{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.AccountService.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}

	return &gober.CreateSessionResponse{
		OfAccount: output.Account,
		Token:     output.Token,
	}, nil
}

func (h *AccountHandler) GetAccount(ctx context.Context, req *gober.GetAccountRequest) (*gober.GetAccountResponse, error) {
	account, err := h.AccountService.GetAccountByID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &gober.GetAccountResponse{
		OfAccount: account,
	}, nil
}

func NewAccountHandler(accountService service.AccountService) (gober.GoberServiceServer, error) {
	return &AccountHandler{
		AccountService: accountService,
	}, nil
}
