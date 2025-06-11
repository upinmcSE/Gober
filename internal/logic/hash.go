package logic

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hashImpl struct{}

// Hash implements Hash.
func (h *hashImpl) Hash(ctx context.Context, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Lỗi khi băm dữ liệu: %v\n", err)
		return "", err
	}

	return string(hashed), nil
}

// IsHashEqual implements Hash.
func (h *hashImpl) IsHashEqual(ctx context.Context, data string, hashed string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(data))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		fmt.Printf("Lỗi khi so sánh băm: %v\n", err)
		return false, err
	}

	return true, nil
}

func NewHash() Hash {
	return &hashImpl{}
}
