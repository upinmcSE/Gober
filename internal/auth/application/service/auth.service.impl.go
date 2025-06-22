package service

import (
	appDto "Gober/internal/auth/application/dto"
	"Gober/internal/auth/domain/model"
	authRepo "Gober/internal/auth/domain/repository"
	hlerDto "Gober/internal/auth/handler/dto"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	authRepo authRepo.AuthRepository
	hasdService HashService
	tokenService TokenService
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


	if account == nil {
		return hlerDto.AccountLoginRes{}, fmt.Errorf("email not found")
	}

	// 2. compare pass
	isEqual, err := a.hasdService.IsHashEqual(ctx, login.Password, account.Password)
	if err != nil{
		return hlerDto.AccountLoginRes{}, fmt.Errorf("failed to compare password: %w", err)
	}

	if !isEqual {
		return hlerDto.AccountLoginRes{}, fmt.Errorf("password is not match")
	}

	// 3. add at vs rt
	claimsAT := jwt.MapClaims{
		"id":   account.ID,
		"role": account.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 1 day
	}

	claimsRT := jwt.MapClaims{
		"id":   account.ID,
		"role": account.Role,
		"exp":  time.Now().Add(time.Hour * 720).Unix(), // 30 days
	}

	accessToken, err := a.tokenService.GenerateToken(claimsAT)
	if err != nil {
		return hlerDto.AccountLoginRes{}, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := a.tokenService.GenerateToken(claimsRT)
	if err != nil {
		return hlerDto.AccountLoginRes{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// 4. return
	return hlerDto.AccountLoginRes{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Id:    account.ID,
		Email: account.Email,
		Role:  string(account.Role),
	}, nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository, 
	hashashService HashService,
	tokenService TokenService,
	) AuthService {
	return &authService{
		authRepo: authRepo,
		hasdService: hashashService,
		tokenService: tokenService,
	}
}
