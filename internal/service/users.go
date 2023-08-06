package service

import (
	"context"
	"fmt"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/internal/repository"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/begenov/region-llc-task/pkg/hash"
	"github.com/begenov/region-llc-task/pkg/logger"
)

type UserService struct {
	userRepo        repository.Users
	hash            hash.PasswordHasher
	manager         *auth.Manager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(userRepo repository.Users, hash hash.PasswordHasher, manager *auth.Manager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *UserService {
	return &UserService{
		userRepo:        userRepo,
		hash:            hash,
		manager:         manager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp domain.UserRequest) (domain.User, error) {
	passwordHash, err := s.hash.GenerateFromPassword(inp.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("s.hash.GenerateFromPassword(): %v", err)
	}

	user := domain.User{
		UserName: inp.UserName,
		Email:    inp.Email,
		Password: passwordHash,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Errorf("s.userRepo.Create(): %v", err)
		return user, fmt.Errorf("s.userRepo.Create(): %v", err)
	}

	return user, nil
}
