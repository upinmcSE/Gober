package service

import (
	appDto "Gober/internal/auth/application/dto"
	"Gober/internal/auth/domain/model"
	authRepo "Gober/internal/auth/domain/repository"
	userRepo "Gober/internal/user/domain/repository"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	authRepo authRepo.AuthRepository
	userRepo userRepo.UserRepository
	hasdService HashService
	tokenService TokenService
}

// Create implements AuthService.
func (a *authService) Create(ctx context.Context, account appDto.AccountAppDTO) (uint, error) {
	// 1. check email in dbs
	existingAccount, err := a.userRepo.EmailExists(ctx, account.Email)
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
func (a *authService) Login(ctx context.Context, login appDto.AccountAppDTO) (appDto.AccountAppLoginDTO, error) {
	// 1. check email in dbs
	account, err := a.userRepo.EmailExists(ctx, login.Email)
	if err != nil {
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("failed to check email: %w", err)
	}


	if account == nil {
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("email not found")
	}

	// 2. compare pass
	isEqual, err := a.hasdService.IsHashEqual(ctx, login.Password, account.Password)
	if err != nil{
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("failed to compare password: %w", err)
	}

	if !isEqual {
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("invalid password")
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
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := a.tokenService.GenerateToken(claimsRT)
	if err != nil {
		return appDto.AccountAppLoginDTO{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// 4. return
	return appDto.AccountAppLoginDTO{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Id:           account.ID,
		Email:        account.Email,
		Role:         string(account.Role),
	}, nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository, 
	userRepo userRepo.UserRepository,
	hashashService HashService,
	tokenService TokenService,
	) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
		hasdService: hashashService,
		tokenService: tokenService,
	}
}
