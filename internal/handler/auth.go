package handler

import (
	"Gober/common/response"
	"Gober/internal/logic"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authLogic logic.AuthLogic
}

func NewAuthHandler(authLogic logic.AuthLogic) *AuthHandler {
	return &AuthHandler{authLogic: authLogic}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// Khai báo request
	var req logic.RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeInvalidParams)
		return
	}

	// Gọi logic đăng ký
	resp, err := h.authLogic.Register(c.Request.Context(), &req)
	if err != nil {
		// Kiểm tra nếu lỗi là CustomError
		if customErr, ok := err.(*response.CustomError); ok {
			response.ErrorResponse(c, customErr.Code)
		} else {
			response.ErrorResponse(c, response.ErrCodeInternalServerError)
		}
		return
	}

	// Trả về phản hồi thành công
	response.SuccessResponse(c, response.ErrCodeSuccess, resp)
}