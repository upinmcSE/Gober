package response

import "fmt"

const (
	ErrCodeSuccess             = 2000
	ErrCodeFail                = 5000
	ErrCodeNotFound            = 4004
	ErrCodeUnauthorized        = 4001
	ErrCodeInvalidParams       = 4002
	ErrCodeNotMatched          = 4003
	ErrCodeEmailAlreadyExists  = 4005
	ErrCodeInternalServerError = 5001
)

var message = map[int]string{
	ErrCodeSuccess:             "success",
	ErrCodeFail:                "fail",
	ErrCodeNotFound:            "not found",
	ErrCodeUnauthorized:        "unauthorized",
	ErrCodeInvalidParams:       "invalid params",
	ErrCodeInternalServerError: "internal server error",
	ErrCodeNotMatched:          "not matched",
	ErrCodeEmailAlreadyExists:  "email already exists",
}

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error: %s (code: %d)", e.Message, e.Code)
}

// NewCustomError tạo một CustomError mới
func NewCustomError(code int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message[code],
	}
}
