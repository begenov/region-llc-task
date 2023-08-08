package mongo

import (
	"testing"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createUser(t *testing.T) domain.User {
	user := domain.User{
		UserName: utils.RandomString(7),
		Email:    utils.RandomEmail(),
		Password: utils.RandomString(7),
		CreateAt: time.Now().Format("2006-01-02"),
	}
	newUser, err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	require.Equal(t, newUser.UserName, user.UserName)
	require.Equal(t, newUser.Password, user.Password)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.CreateAt, user.CreateAt)

	require.NotEmpty(t, newUser.ID)

	return newUser
}

func TestUserRepo_Create(t *testing.T) {
	createUser(t)
}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	user := createUser(t)
	require.NotEmpty(t, user)

	userE, err := userRepo.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userE)

	require.Equal(t, user.ID, userE.ID)
	require.Equal(t, user.Email, userE.Email)
	require.Equal(t, user.UserName, userE.UserName)
	require.Equal(t, user.CreateAt, userE.CreateAt)
	require.Equal(t, user.Password, userE.Password)

	_, err = userRepo.GetUserByEmail(ctx, "no-email")
	require.Equal(t, err, domain.ErrNotFound)
}

func TestUserRepo_GetUserByID(t *testing.T) {
	user := createUser(t)
	require.NotEmpty(t, user)

	userI, err := userRepo.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userI)

	require.Equal(t, user.ID, userI.ID)
	require.Equal(t, user.Email, userI.Email)
	require.Equal(t, user.UserName, userI.UserName)
	require.Equal(t, user.CreateAt, userI.CreateAt)
	require.Equal(t, user.Password, userI.Password)

	_, err = userRepo.GetUserByID(ctx, primitive.ObjectID{})
	require.Equal(t, err, domain.ErrNotFound)
}

func setSession(t *testing.T) (domain.User, string) {
	session := domain.Session{
		RefreshToken: utils.RandomString(10),
		ExpirationAt: time.Now().Add(time.Minute),
	}

	user := createUser(t)
	require.NotEmpty(t, user)

	err := userRepo.SetSession(ctx, user.ID, session)
	require.NoError(t, err)

	return user, session.RefreshToken
}

func TestUserRepo_SetSession(t *testing.T) {
	setSession(t)
}

func TestUserRepo_GetByRefreshToken(t *testing.T) {
	user, refreshToken := setSession(t)
	require.NotEmpty(t, user)
	require.NotEmpty(t, refreshToken)

	userR, err := userRepo.GetByRefreshToken(ctx, refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, userR)

	require.Equal(t, user.ID, userR.ID)
	require.Equal(t, user.UserName, userR.UserName)
	require.Equal(t, user.Password, userR.Password)
	require.Equal(t, user.Email, userR.Email)
	require.Equal(t, user.CreateAt, userR.CreateAt)

	require.NotEmpty(t, userR.Session)

	_, err = userRepo.GetByRefreshToken(ctx, "")
	require.Error(t, err)
}
