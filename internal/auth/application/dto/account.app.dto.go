package dto

type AccountAppDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountAppLoginDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Id           uint   `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}