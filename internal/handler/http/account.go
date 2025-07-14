package http

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AccountHandler struct {
	client gober.GoberServiceClient
}

func (ah *AccountHandler) CreateSessionHandler(c *gin.Context) {
	var req gober.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	resp, err := ah.client.CreateSession(c, &req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Session created successfully", resp)

}

func (ah *AccountHandler) CreateHandler(c *gin.Context) {

	var req gober.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	resp, err := ah.client.CreateAccount(c, &req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 201, "Account created successfully", resp)

}

func (ah *AccountHandler) GetAccountHandler(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Account ID is required",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	// Convert accountID to int64
	id, err := strconv.ParseUint(accountID, 10, 64)

	resp, err := ah.client.GetAccount(c, &gober.GetAccountRequest{AccountId: id})
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Account retrieved", resp)
}

func NewAccountHandler(client gober.GoberServiceClient) *AccountHandler {
	return &AccountHandler{
		client: client,
	}
}
