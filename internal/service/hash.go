package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Hash interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hash struct {
}

func NewHash() Hash {
	return &hash{}
}

func (h hash) Hash(_ context.Context, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to hash data")
	}

	return string(hashed), nil
}

func (h hash) IsHashEqual(_ context.Context, data string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(data)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, status.Error(codes.Internal, "failed to check if data equal hash")
	}

	return true, nil
}