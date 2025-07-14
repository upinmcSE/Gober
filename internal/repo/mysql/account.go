package mysql

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRole string

const (
	Manager  UserRole = "manager"
	Attendee UserRole = "attendee"
)

type Account struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" gorm:"text;not null"`
	Role      UserRole  `json:"role" gorm:"text;default:attendee"`
	Password  string    `json:"-"` // Do not compute the password in json
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	IsDeleted bool      `json:"isDeleted" gorm:"bool;default:false;column:is_deleted"`
}

type AccountDatabase interface {
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccountByID(ctx context.Context, id uint64) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
}

type accountDatabase struct {
	db *gorm.DB
}

// CreateAccount implements AccountDatabase.
func (a *accountDatabase) CreateAccount(ctx context.Context, account *Account) (*Account, error) {
	var result *Account
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(account).Error; err != nil {
			return err // rollback
		}
		// Nếu muốn thao tác thêm ở đây, bạn dùng `tx` tiếp tục

		// Gán kết quả để return sau khi commit
		result = account
		return nil // commit
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get last inserted id")
	}
	return result, nil
}

// GetAccountByEmail implements AccountDatabase.
func (a *accountDatabase) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var account *Account
	err := a.db.WithContext(ctx).
		Where(&Account{Email: email, IsDeleted: false}).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // email chưa tồn tại
		}
		return nil, err
	}

	return account, nil
}

// GetAccountByID implements AccountDatabase.
func (a *accountDatabase) GetAccountByID(ctx context.Context, id uint64) (*Account, error) {
	var account Account
	err := a.db.WithContext(ctx).
		Where(&Account{ID: id, IsDeleted: false}).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func NewAccountDatabase(db *gorm.DB) AccountDatabase {
	return &accountDatabase{db: db}
}
