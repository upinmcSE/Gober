package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type HashService interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hashService struct{}

// Hash implements HashService.
func (h *hashService) Hash(ctx context.Context, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// IsHashEqual implements HashService.
func (h *hashService) IsHashEqual(ctx context.Context, data string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(data)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, fmt.Errorf("password does not match: %w", err) 
		}
		return false, fmt.Errorf("error comparing hash: %w", err)
	}
	return true, nil
}

func NewHashService() HashService {
	return &hashService{}
}
