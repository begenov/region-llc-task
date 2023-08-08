package service

import (
	"context"
	"net/mail"
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
	redisRepo       repository.Redis
}

func NewUserService(userRepo repository.Users, hash hash.PasswordHasher, manager *auth.Manager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, redisRepo repository.Redis) *UserService {
	return &UserService{
		userRepo:        userRepo,
		hash:            hash,
		manager:         manager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		redisRepo:       redisRepo,
	}
}

func (s *UserService) SignUp(ctx context.Context, inp domain.UserRequest) (domain.User, error) {

	if err := validateUser(inp.Email, inp.Password); err != nil {
		logger.Errorf("validateUser(): %v", err)
		return domain.User{}, err
	}

	if inp.UserName == "" && len(inp.UserName) < 4 {
		return domain.User{}, domain.ErrIncorrectUserName
	}

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
		CreateAt: time.Now().Format("2006-01-02"),
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Errorf("s.userRepo.Create(): %v", err)
		return user, err
	}

	return user, nil
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (domain.Token, error) {

	if err := validateUser(email, password); err != nil {
		logger.Errorf("validateUser(): %v", err)
		return domain.Token{}, err
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Errorf("s.userRepo.GetUserByEmail(): %v", err)
		return domain.Token{}, err
	}

	if err := s.hash.CompareHashAndPassword(user.Password, password); err != nil {
		logger.Errorf("s.hash.CompareHashAndPassword(): %v", err)
		return domain.Token{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (domain.Token, error) {
	idString, err := s.redisRepo.Get(refreshToken)
	if err == nil {
		idObject, err := primitive.ObjectIDFromHex(idString)
		if err == nil {
			return s.createSession(ctx, idObject)
		}
	}

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

	id := userId.Hex()

	res.AccessToken, err = s.manager.NewJWT(id, s.accessTokenTTL)
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

	err = s.redisRepo.Set(session.RefreshToken, id, s.accessTokenTTL*2)
	if err != nil {
		logger.Errorf("s.redisRepo.Set(): %v", err)
		return res, err
	}

	err = s.userRepo.SetSession(ctx, userId, session)

	return res, err
}

func validateUser(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain.ErrIncorrectEmailAddress
	}

	if len(password) < 6 {
		return domain.ErrIncorrectPassword
	}

	return nil
}
