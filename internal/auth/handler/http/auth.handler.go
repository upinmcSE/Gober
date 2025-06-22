package http

import (
	appDto "Gober/internal/auth/application/dto"
	"Gober/internal/auth/application/service"
	hlerDto "Gober/internal/auth/handler/dto"
	"Gober/pkg/response"
	"Gober/utils"
	"fmt"
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
	fmt.Println("---> RegisterUser")

	var req hlerDto.AccountRegisterReq
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

	// Create account -> dto application
	account := appDto.AccountAppDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	accountId, err := ah.service.Create(ctx, account)
	if err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Registration failed", err.Error())
	}

	return accountId, nil
}

func (ah *AuthHandler) LoginUser(ctx *gin.Context) (res interface{}, err error) {
	fmt.Println("---> LoginUser")

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

	accountLoginRes, err := ah.service.Login(ctx, req)
	if err != nil {
		return nil, response.NewAPIError(http.StatusOK, "Login failed", err.Error())
	}

	return accountLoginRes, nil
}