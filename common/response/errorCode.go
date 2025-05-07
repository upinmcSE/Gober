package response

const (
	ErrCodeSuccess       = 2000
	ErrCodeFail          = 5000
	ErrCodeNotFound      = 4004
	ErrCodeUnauthorized  = 4001
	ErrCodeInvalidParams = 4002
)

var message = map[int]string{
	ErrCodeSuccess:       "success",
	ErrCodeFail:          "fail",
	ErrCodeNotFound:      "not found",
	ErrCodeUnauthorized:  "unauthorized",
	ErrCodeInvalidParams: "invalid params",
}
