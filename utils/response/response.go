package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeInternal     ErrorCode = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

type APIResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (ae *AppError) Error() string {
	return ""
}

func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	resp := APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	ctx.JSON(status, resp)
}

func ResponseError(ctx *gin.Context, status int, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

func HandleGrpcError(c *gin.Context, err error) {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		ResponseError(c, http.StatusInternalServerError, &AppError{
			Message: "Lỗi không xác định",
			Code:    ErrCodeInternal,
		})
		return
	}

	switch grpcStatus.Code() {
	case codes.NotFound:
		ResponseError(c, http.StatusNotFound, &AppError{
			Message: grpcStatus.Message(),
			Code:    ErrCodeNotFound,
		})
	case codes.AlreadyExists:
		ResponseError(c, http.StatusConflict, &AppError{
			Message: grpcStatus.Message(),
			Code:    ErrCodeConflict,
		})
	case codes.Unauthenticated:
		ResponseError(c, http.StatusUnauthorized, &AppError{
			Message: grpcStatus.Message(),
			Code:    ErrCodeUnauthorized,
		})
	case codes.InvalidArgument:
		ResponseError(c, http.StatusBadRequest, &AppError{
			Message: grpcStatus.Message(),
			Code:    ErrCodeBadRequest,
		})
	case codes.Internal:
		ResponseError(c, http.StatusInternalServerError, &AppError{
			Message: "Lỗi máy chủ nội bộ",
			Code:    ErrCodeInternal,
		})
	default:
		ResponseError(c, http.StatusInternalServerError, &AppError{
			Message: grpcStatus.Message(),
			Code:    ErrCodeInternal,
		})
	}
}
