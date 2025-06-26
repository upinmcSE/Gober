package dto

type UserInfo struct {
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AccountRegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type AccountRegisterRes struct {
	AccountID uint   `json:"account_id"`
	Message   string `json:"message"`
}

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type AccountLoginRes struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	User         UserInfo `json:"user"`
}