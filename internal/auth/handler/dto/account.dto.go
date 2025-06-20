package dto

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type AccountLoginRes struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Id           uint   `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
}