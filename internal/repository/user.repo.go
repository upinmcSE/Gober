package repository

import "fmt"

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) GetUserById(id string) string {
	return "user" + fmt.Sprint(id)
}