package http

import (
	appDto "Gober/internal/auth/application/dto"
	"Gober/internal/auth/application/service"
	hlerDto "Gober/internal/auth/handler/dto"
	"Gober/pkg/response"
	"Gober/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}


func (ah *AuthHandler) RegisterUser(ctx *gin.Context) (res interface{}, err error) {
	var req hlerDto.AccountRegisterReq

	// 1: Parse JSON → struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Invalid request", err.Error())
	}

	// 2: Validate business rules
	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusOK, "Invalid request", "Validation not found in context")
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// 3: Create account -> dto application
	account := appDto.AccountAppDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	// 4: Gọi service to create account
	accountId, err := ah.service.Create(ctx, account)
	if err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Registration failed", err.Error())
	}

	// 5: Convert to handler dto
	handlerResult := hlerDto.AccountRegisterRes{
		AccountID: accountId,
		Message:   "Account created successfully",
	}

	return handlerResult, nil
}

func (ah *AuthHandler) LoginUser(ctx *gin.Context) (res interface{}, err error) {
	var req hlerDto.AccountLoginReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Invalid request", err.Error())
	}

	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusOK, "Invalid request", "Validation not found in context")
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// Convert hlerDto to appDto
	accountAppDTO := appDto.AccountAppDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	serviceResult, err := ah.service.Login(ctx, accountAppDTO)
	if err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Login failed", err.Error())
	}
	// Convert appDto to hlerDto

	userInfo := hlerDto.UserInfo{
		Id:        serviceResult.Id,
		Email:     serviceResult.Email,
		Role:      serviceResult.Role,
		CreatedAt: serviceResult.CreatedAt,
		UpdatedAt: serviceResult.UpdatedAt,
	}

	handlerResult := hlerDto.AccountLoginRes{
		Token:        serviceResult.Token,
		RefreshToken: serviceResult.RefreshToken,
		User:         userInfo,
	}

	return handlerResult, nil
}