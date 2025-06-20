package service

import (
	appDto "Gober/internal/auth/application/dto"
	"Gober/internal/auth/domain/model"
	authRepo "Gober/internal/auth/domain/repository"
	hlerDto "Gober/internal/auth/handler/dto"
	"context"
	"fmt"
)

type authService struct {
	authRepo authRepo.AuthRepository
	hasdService HashService
}

// Create implements AuthService.
func (a *authService) Create(ctx context.Context, account appDto.AccountAppDTO) (uint, error) {
	// 1. check email in dbs
	existingAccount, err := a.authRepo.EmailExists(ctx, account.Email)
	if err != nil {
		return 0, err
	}
	if existingAccount != nil {
		return 0, fmt.Errorf("email already exists")
	}

	// 2. hash pass
	hashedPassword, err := a.hasdService.Hash(ctx, account.Password)
	if err != nil {
		return 0, err
	}

	// 3. create account
	newAccount := &model.Account{
		Email:    account.Email,
		Password: hashedPassword,
	}

	accountCreated, err := a.authRepo.CreateUser(ctx, newAccount)
	if err != nil {
		return 0, err
	}

	return accountCreated.ID, nil
}

// Login implements AuthService.
func (a *authService) Login(ctx context.Context, login hlerDto.AccountLoginReq) (hlerDto.AccountLoginRes, error) {
	// 1. check email in dbs
	account, err := a.authRepo.EmailExists(ctx, login.Email)
	if err != nil {
		return hlerDto.AccountLoginRes{}, err
	}

	// 2. compare pass
	isEqual, err := a.hasdService.IsHashEqual(ctx, login.Password, account.Password)
	if err != nil{
		return hlerDto.AccountLoginRes{}, err
	}

	if !isEqual {
		return hlerDto.AccountLoginRes{}, fmt.Errorf("password is not match")
	}

	// 3. add at vs rt

	return hlerDto.AccountLoginRes{
		Token:        "",
		RefreshToken: "",
		Id:    account.ID,
		Email: account.Email,
		Role:  string(account.Role),
	}, nil
}

func NewAuthService(authRepo authRepo.AuthRepository, hashashService HashService) AuthService {
	return &authService{
		authRepo: authRepo,
		hasdService: hashashService,
	}
}
