package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	Manager  UserRole = "manager"
	attendee UserRole = "attendee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" gorm:"text;not null"`
	Name	  string    `json:"name" gorm:"text;not null"`
	Phone     string    `json:"phone" gorm:"text;not null"`
	Role      UserRole  `json:"role" gorm:"text;default:attendee"`
	Password  string    `json:"-"` // Do not compute the password in json
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserDataAccessor interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UserExist(ctx context.Context, email string) (bool, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserList(ctx context.Context, offset, limit int) ([]*User, error)
	GetUserCount(ctx context.Context) (int64, error)
}

type userDataAccessorImpl struct {
	db *gorm.DB
}

// CreateUser implements UserDataAccessor.
func (u *userDataAccessorImpl) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return gorm.ErrInvalidData
	}
	
	// Set default role if not provided
	if user.Role == "" {
		user.Role = attendee
	}

	// Create the user in the database
	return u.db.WithContext(ctx).Create(user).Error
}


// DeleteUser implements UserDataAccessor.
func (u *userDataAccessorImpl) DeleteUser(ctx context.Context, id uint) error {
	panic("unimplemented")
}

// GetUserByEmail implements UserDataAccessor.
func (u *userDataAccessorImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err 
	}
	return &user, nil
}

// UserExist implements UserDataAccessor.
func (u *userDataAccessorImpl) UserExist(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserByID implements UserDataAccessor.
func (u *userDataAccessorImpl) GetUserByID(ctx context.Context, id uint) (*User, error) {
	panic("unimplemented")
}

// GetUserCount implements UserDataAccessor.
func (u *userDataAccessorImpl) GetUserCount(ctx context.Context) (int64, error) {
	panic("unimplemented")
}

// GetUserList implements UserDataAccessor.
func (u *userDataAccessorImpl) GetUserList(ctx context.Context, offset int, limit int) ([]*User, error) {
	panic("unimplemented")
}

// UpdateUser implements UserDataAccessor.
func (u *userDataAccessorImpl) UpdateUser(ctx context.Context, user *User) error {
	panic("unimplemented")
}


func NewUserDataAccessor(db *gorm.DB) UserDataAccessor {
	return &userDataAccessorImpl{db: db}
}
