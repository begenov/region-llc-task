package service

import (
	"context"
	"testing"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	mocksRepo "github.com/begenov/region-llc-task/internal/repository/mocks"
	"github.com/begenov/region-llc-task/pkg/auth"
	mocksHash "github.com/begenov/region-llc-task/pkg/hash/mocks"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

func TestUserService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksRepo.NewMockUsers(ctrl)
	redisRepo := mocksRepo.NewMockRedis(ctrl)
	hash := mocksHash.NewMockPasswordHasher(ctrl)

	manager, err := auth.NewManager("qwerty")
	if err != nil {
		t.Fatal(err)
	}

	userService := *NewUserService(userRepo, hash, manager, time.Minute, time.Minute, redisRepo)

	type args struct {
		ctx context.Context
		inp domain.UserRequest
	}

	tests := []struct {
		name          string
		args          args
		buildStubs    func(user domain.UserRequest)
		checkResponse func(user domain.User, req domain.UserRequest, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(7),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, domain.ErrNotFound)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(1).Return(user.Password, nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
					CreateAt: time.Now().Format("2006-01-02"),
				}, nil)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.NoError(t, err)
				require.Equal(t, req.UserName, user.UserName)
				require.Equal(t, req.Email, user.Email)
			},
		},
		{
			name: "Incorrect Email Address",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomString(10),
					Password: utils.RandomString(7),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrIncorrectEmailAddress)
			},
		},
		{
			name: "Incorrect Password",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: "",
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrIncorrectPassword)
			},
		},
		{
			name: "Incorrect UserName",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: "",
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrIncorrectUserName)
			},
		},
		{
			name: "Email Already Exists",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
				}, nil)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrEmailAlreadyExists)
			},
		},
		{
			name: "Generate Password Hash",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, domain.ErrNotFound)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(1).Return("", domain.ErrInternalServer)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
		{
			name: "Generate Password Hash",
			args: args{
				ctx: ctx,
				inp: domain.UserRequest{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.UserRequest) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, domain.ErrNotFound)
				hash.EXPECT().GenerateFromPassword(gomock.Any()).Times(1).Return(user.Password, nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, domain.ErrInternalServer)
			},
			checkResponse: func(user domain.User, req domain.UserRequest, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.buildStubs(tt.args.inp)

			user, err := userService.SignUp(tt.args.ctx, tt.args.inp)

			tt.checkResponse(user, tt.args.inp, err)
		})
	}
}

func TestUserService_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksRepo.NewMockUsers(ctrl)
	redisRepo := mocksRepo.NewMockRedis(ctrl)
	hash := mocksHash.NewMockPasswordHasher(ctrl)

	manager, err := auth.NewManager("qwerty")
	if err != nil {
		t.Fatal(err)
	}
	userService := *NewUserService(userRepo, hash, manager, time.Minute, time.Minute, redisRepo)

	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	tests := []struct {
		name          string
		args          args
		buildStubs    func(email, password string)
		checkResponse func(token domain.Token, err error)
	}{
		{
			name: "OK",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: utils.RandomString(10),
			},
			buildStubs: func(email, password string) {
				id := primitive.NewObjectID()
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(1).Return(domain.User{
					ID:       id,
					UserName: utils.RandomString(10),
					Email:    email,
					Password: password,
					CreateAt: time.Now().Format("2006-01-02"),
				}, nil)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(token domain.Token, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, token.AccessToken)
				require.NotEmpty(t, token.RefreshToken)
			},
		},
		{
			name: "Incorrect Email Address",
			args: args{
				ctx:      ctx,
				email:    "",
				password: utils.RandomString(10),
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(0)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(0)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrIncorrectEmailAddress)
			},
		},
		{
			name: "Incorrect password",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: "",
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(0)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(0)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrIncorrectPassword)
			},
		},
		{
			name: "not found user",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: utils.RandomString(10),
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(1).Return(domain.User{}, domain.ErrNotFound)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(0)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
		{
			name: "internal server",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: utils.RandomString(10),
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(1).Return(domain.User{}, domain.ErrInternalServer)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(0)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
		{
			name: "incorect password",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: utils.RandomString(10),
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(1).Return(domain.User{
					ID:       primitive.NewObjectID(),
					UserName: utils.RandomString(10),
					Email:    email,
					Password: password,
					CreateAt: time.Now().Format("2006-01-02"),
				}, nil)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(1).Return(domain.ErrInternalServer)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
		{
			name: "set redis",
			args: args{
				ctx:      ctx,
				email:    utils.RandomEmail(),
				password: utils.RandomString(10),
			},

			buildStubs: func(email, password string) {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Times(1).Return(domain.User{
					ID:       primitive.NewObjectID(),
					UserName: utils.RandomString(10),
					Email:    email,
					Password: password,
					CreateAt: time.Now().Format("2006-01-02"),
				}, nil)
				hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(domain.ErrInternalServer)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(tt.args.email, tt.args.password)

			token, err := userService.SignIn(tt.args.ctx, tt.args.email, tt.args.password)

			tt.checkResponse(token, err)
		})
	}
}

func TestUserService_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksRepo.NewMockUsers(ctrl)
	redisRepo := mocksRepo.NewMockRedis(ctrl)
	hash := mocksHash.NewMockPasswordHasher(ctrl)

	manager, err := auth.NewManager("qwerty")
	if err != nil {
		t.Fatal(err)
	}

	userService := *NewUserService(userRepo, hash, manager, time.Minute, time.Minute, redisRepo)

	type args struct {
		ctx     context.Context
		refresh string
	}

	tests := []struct {
		name          string
		args          args
		buildStubs    func(refresh string)
		checkResponse func(token domain.Token, err error)
	}{
		{
			name: "OK-1",
			args: args{
				ctx:     ctx,
				refresh: utils.RandomString(10),
			},
			buildStubs: func(refresh string) {
				id := primitive.NewObjectID()
				redisRepo.EXPECT().Get(gomock.Any()).Times(1).Return(id.Hex(), nil)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				userRepo.EXPECT().GetByRefreshToken(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(token domain.Token, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, token.AccessToken)
				require.NotEmpty(t, token.RefreshToken)

			},
		},
		{
			name: "OK-2",
			args: args{
				ctx:     ctx,
				refresh: utils.RandomString(10),
			},
			buildStubs: func(refresh string) {
				id := primitive.NewObjectID()
				redisRepo.EXPECT().Get(gomock.Any()).Times(1).Return("", nil)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				userRepo.EXPECT().GetByRefreshToken(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{
					ID:       id,
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				}, nil)
			},
			checkResponse: func(token domain.Token, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, token.AccessToken)
				require.NotEmpty(t, token.RefreshToken)

			},
		},
		{
			name: "error not found refresh token",
			args: args{
				ctx:     ctx,
				refresh: utils.RandomString(10),
			},
			buildStubs: func(refresh string) {
				redisRepo.EXPECT().Get(gomock.Any()).Times(1).Return("", nil)
				redisRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().SetSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().GetByRefreshToken(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, domain.ErrNotFound)
			},
			checkResponse: func(token domain.Token, err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(tt.args.refresh)

			token, err := userService.RefreshTokens(tt.args.ctx, tt.args.refresh)

			tt.checkResponse(token, err)
		})
	}
}
