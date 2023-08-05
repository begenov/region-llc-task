package service

import "github.com/begenov/region-llc-task/internal/repository"

type UserService struct {
	userRepo repository.Users
}

func NewUserService(userRepo repository.Users) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
