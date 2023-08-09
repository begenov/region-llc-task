package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/begenov/region-llc-task/internal/domain"
	serviceMocks "github.com/begenov/region-llc-task/internal/service/mocks"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestServer_userSignUp(t *testing.T) {
	tests := []struct {
		name          string
		inp           domain.UserRequest
		buildStubs    func(service *serviceMocks.MockUsers, inp domain.UserRequest)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "Status ok",
			inp: domain.UserRequest{
				UserName: utils.RandomString(10),
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(10),
			},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserRequest) {
				service.EXPECT().SignUp(gomock.Any(), inp).Times(1).Return(domain.User{
					ID:       primitive.NewObjectID(),
					UserName: inp.UserName,
					Email:    inp.Email,
					Password: inp.Password,
				}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "Status Bad Request",
			inp:  domain.UserRequest{},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserRequest) {
				service.EXPECT().SignUp(gomock.Any(), inp).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "Internal Server",
			inp: domain.UserRequest{
				UserName: utils.RandomString(10),
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(10),
			},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserRequest) {
				service.EXPECT().SignUp(gomock.Any(), inp).Times(1).Return(domain.User{}, domain.ErrInternalServer)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := serviceMocks.NewMockUsers(ctrl)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			handler := &Server{
				userService: userService,
			}
			tt.buildStubs(userService, tt.inp)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/sign-up", server.URL)
			body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_userSignIn(t *testing.T) {

	tests := []struct {
		name          string
		inp           domain.UserSignInRequest
		buildStubs    func(service *serviceMocks.MockUsers, inp domain.UserSignInRequest)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "Status OK",
			inp: domain.UserSignInRequest{
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(10),
			},

			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserSignInRequest) {
				service.EXPECT().SignIn(gomock.Any(), inp.Email, inp.Password).Times(1).Return(domain.Token{
					RefreshToken: utils.RandomString(10),
					AccessToken:  utils.RandomString(25),
				}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "bad request",
			inp:  domain.UserSignInRequest{},

			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserSignInRequest) {
				service.EXPECT().SignIn(gomock.Any(), inp.Email, inp.Password).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "internal server",
			inp: domain.UserSignInRequest{
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(10),
			},

			buildStubs: func(service *serviceMocks.MockUsers, inp domain.UserSignInRequest) {
				service.EXPECT().SignIn(gomock.Any(), inp.Email, inp.Password).Times(1).Return(domain.Token{}, domain.ErrInternalServer)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := serviceMocks.NewMockUsers(ctrl)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			handler := &Server{
				userService: userService,
			}
			tt.buildStubs(userService, tt.inp)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/sign-in", server.URL)
			body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_userRefresh(t *testing.T) {

	tests := []struct {
		name          string
		inp           domain.RefreshToken
		buildStubs    func(service *serviceMocks.MockUsers, inp domain.RefreshToken)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "Status Ok",
			inp: domain.RefreshToken{
				RefreshToken: utils.RandomString(15),
			},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.RefreshToken) {
				service.EXPECT().RefreshTokens(gomock.Any(), inp.RefreshToken).Times(1).Return(domain.Token{
					AccessToken:  utils.RandomString(10),
					RefreshToken: utils.RandomString(10),
				}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "bad request",
			inp:  domain.RefreshToken{},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.RefreshToken) {
				service.EXPECT().RefreshTokens(gomock.Any(), inp.RefreshToken).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "internal server",
			inp: domain.RefreshToken{
				RefreshToken: utils.RandomString(10),
			},
			buildStubs: func(service *serviceMocks.MockUsers, inp domain.RefreshToken) {
				service.EXPECT().RefreshTokens(gomock.Any(), inp.RefreshToken).Times(1).Return(domain.Token{}, domain.ErrInternalServer)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := serviceMocks.NewMockUsers(ctrl)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			handler := &Server{
				userService: userService,
			}
			tt.buildStubs(userService, tt.inp)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/auth/refresh", server.URL)
			body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
