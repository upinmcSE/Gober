package service

import "Gober/internal/repository"

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepo(),
	}
}

func (us *UserService) GetUserById(id string) string {
	return us.userRepo.GetUserById(id)
}