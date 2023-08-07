package service

import (
	"context"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/internal/repository"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/begenov/region-llc-task/pkg/hash"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	_, err := s.userRepo.GetUserByEmail(ctx, inp.Email)
	if err == nil {
		logger.Errorf("s.registerRepo.GetRegisterUsername(): %v", err)
		return domain.User{}, domain.ErrEmailAlreadyExists
	}

	passwordHash, err := s.hash.GenerateFromPassword(inp.Password)
	if err != nil {
		logger.Errorf("s.hash.GenerateFromPassword(): %v", err)
		return domain.User{}, domain.ErrInternalServer
	}

	user := domain.User{
		UserName: inp.UserName,
		Email:    inp.Email,
		Password: passwordHash,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Errorf("s.userRepo.Create(): %v", err)
		return user, err
	}

	return user, nil
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (domain.Token, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Errorf("s.userRepo.GetUserByEmail(): %v", err)
		return domain.Token{}, err
	}

	logger.Infof("%s\t%s", user.Password, password)
	if err := s.hash.CompareHashAndPassword(user.Password, password); err != nil {
		logger.Errorf("s.hash.CompareHashAndPassword(): %v", err)
		return domain.Token{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (domain.Token, error) {
	user, err := s.userRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.Token{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) createSession(ctx context.Context, userId primitive.ObjectID) (domain.Token, error) {
	var (
		res domain.Token
		err error
	)

	res.AccessToken, err = s.manager.NewJWT(userId.Hex(), s.accessTokenTTL)
	if err != nil {
		logger.Errorf("s.manager.NewJWT(): %v", err)
		return res, err
	}

	res.RefreshToken, err = s.manager.NewRefreshToken()
	if err != nil {
		logger.Errorf("s.manager.NewRefreshToken(): %v", err)
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpirationAt: time.Now().Add(s.refreshTokenTTL),
	}

	err = s.userRepo.SetSession(ctx, userId, session)

	return res, err
}
