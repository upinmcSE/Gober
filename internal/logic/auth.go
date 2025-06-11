package logic

import (
	"Gober/common/response"
	"Gober/internal/database"
	"context"
	"regexp"
)

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required,min=2,max=100"`
	Phone      string `json:"phone"`
	Password   string `json:"password" binding:"required,min=6,max=100"`
	RePassword string `json:"rePassword" binding:"required,min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AuthLogic interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*UserResponse, error)
}

type authLogicImpl struct {
	userDataAccessor database.UserDataAccessor
	hashLogic        Hash
}

// validate phone
func validatePhone(phone string) bool {
	pattern := `^0[0-9]{9}$`

	// Biên dịch regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	// Kiểm tra số điện thoại
	return re.MatchString(phone)
}



// GetUserByID implements UserLogic.
func (u *authLogicImpl) GetUserByID(ctx context.Context, id uint) (*UserResponse, error) {
	panic("unimplemented")
}

// Login implements UserLogic.
func (u *authLogicImpl) Login(ctx context.Context, req *LoginRequest) (*UserResponse, error) {
	panic("unimplemented")
}

// Register implements UserLogic.
func (u *authLogicImpl) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// check user exist
	user, err := u.userDataAccessor.UserExist(ctx, req.Email)
	if err != nil {
		return nil, response.NewCustomError(response.ErrCodeInternalServerError)
	}

	if user {
		return nil, response.NewCustomError(response.ErrCodeEmailAlreadyExists)
	}

	// insert user
	hashedPassword, err := u.hashLogic.Hash(ctx, req.Password)
    if err != nil {
        return nil, response.NewCustomError(response.ErrCodeInternalServerError)
    }

    newUser := &database.User{
        Email:    req.Email,
        Name:     req.Name,
        Phone:    req.Phone,
        Password: hashedPassword,
    }

	
	err = u.userDataAccessor.CreateUser(ctx, newUser)
	if err != nil {
		return nil, response.NewCustomError(response.ErrCodeInternalServerError)
	}

	return &UserResponse{
		Email:     newUser.Email,
		Name:      newUser.Name,
		Phone:     newUser.Phone,
		Role:      string(newUser.Role),
		CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, err
	
}

// NewUserLogic creates a new instance of UserLogic.
func NewAuthLogic(userDataAccessor database.UserDataAccessor, hashLogic Hash) AuthLogic {
	return &authLogicImpl{
		userDataAccessor: userDataAccessor,
		hashLogic:        hashLogic,
	}
}
