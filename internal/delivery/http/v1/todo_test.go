package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	repoMocks "github.com/begenov/region-llc-task/internal/repository/mocks"
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func addAuthorization(t *testing.T, request *http.Request, token auth.TokenManager, authorizationType string, username string, duration time.Duration) {
	accessToken, err := token.NewJWT(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestServer_createTodo(t *testing.T) {

	tests := []struct {
		name          string
		inp           domain.TodoRequest
		userID        string
		setupAuth     func(request *http.Request, id string, token auth.TokenManager)
		buildStubs    func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "status ok",
			inp: domain.TodoRequest{
				Title:    utils.RandomString(10),
				ActiveAt: time.Now().Format(domain.Format),
			},
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},

			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(count, nil)
				userRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{
					ID:       primitive.NewObjectID(),
					UserName: utils.RandomString(10),
					Email:    utils.RandomString(10),
					Password: utils.RandomString(10),
				}, nil)

				todoRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    inp.Title,
					Author:   utils.RandomString(10),
					ActiveAt: time.Now().Format(domain.Format),
					Status:   domain.Active,
				}, nil)
			},
			userID: primitive.NewObjectID().Hex(),
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "bad request",
			inp: domain.TodoRequest{
				Title:    "",
				ActiveAt: time.Now().Format(domain.Format),
			},
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest) {
				// var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				todoRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "auth error",
			inp: domain.TodoRequest{
				Title:    "",
				ActiveAt: time.Now().Format(domain.Format),
			},
			userID: utils.RandomString(10),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest) {
				// var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				todoRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
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

			todoRepo := repoMocks.NewMockTodo(ctrl)
			userRepo := repoMocks.NewMockUsers(ctrl)
			todoService := service.NewTodoService(todoRepo, userRepo)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			token, err := auth.NewManager("qwerty")
			require.NoError(t, err)
			handler := &Server{
				todoService:  todoService,
				tokenManager: token,
			}
			tt.buildStubs(userRepo, todoRepo, tt.inp)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/todo-list/todo", server.URL)
			body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			tt.setupAuth(request, tt.userID, token)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_updateTodo(t *testing.T) {

	tests := []struct {
		name          string
		inp           domain.TodoRequest
		uri           domain.TodoURI
		userID        primitive.ObjectID
		todoID        primitive.ObjectID
		setupAuth     func(request *http.Request, id string, token auth.TokenManager)
		buildStubs    func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest, id primitive.ObjectID, userID primitive.ObjectID)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			inp: domain.TodoRequest{
				Title:    utils.RandomString(10),
				ActiveAt: time.Now().Format(domain.Format),
			},
			uri: domain.TodoURI{},
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest, id primitive.ObjectID, userID primitive.ObjectID) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any()).Times(1).Return(domain.Todo{
					ID:       id,
					UserID:   userID,
					Title:    inp.Title,
					ActiveAt: inp.ActiveAt,
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				}, nil)
			},

			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "title empty",
			inp: domain.TodoRequest{
				Title:    "",
				ActiveAt: time.Now().Format(domain.Format),
			},
			uri: domain.TodoURI{},
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, inp domain.TodoRequest, id primitive.ObjectID, userID primitive.ObjectID) {
				// var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any()).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any()).Times(0)
			},

			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepo := repoMocks.NewMockTodo(ctrl)
			userRepo := repoMocks.NewMockUsers(ctrl)
			todoService := service.NewTodoService(todoRepo, userRepo)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			token, err := auth.NewManager("qwerty")
			require.NoError(t, err)
			handler := &Server{
				todoService:  todoService,
				tokenManager: token,
			}
			tt.buildStubs(userRepo, todoRepo, tt.inp, tt.todoID, tt.userID)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/todo-list/todo/%s", server.URL, tt.todoID.Hex())
			body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
			tt.setupAuth(request, tt.userID.Hex(), token)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_deleteTodo(t *testing.T) {

	tests := []struct {
		name          string
		uri           domain.TodoURI
		userID        string
		setupAuth     func(request *http.Request, id string, token auth.TokenManager)
		buildStubs    func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "status ok",
			uri: domain.TodoURI{
				ID: primitive.NewObjectID().Hex(),
			},
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo) {
				todoRepo.EXPECT().DeleteTodoByID(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		// TODO: Add test cases.
		{
			name: "todo invalid id",
			uri: domain.TodoURI{
				ID: "primitive.NewObjectID().Hex()",
			},
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo) {
				todoRepo.EXPECT().DeleteTodoByID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepo := repoMocks.NewMockTodo(ctrl)
			userRepo := repoMocks.NewMockUsers(ctrl)
			todoService := service.NewTodoService(todoRepo, userRepo)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			token, err := auth.NewManager("qwerty")
			require.NoError(t, err)
			handler := &Server{
				todoService:  todoService,
				tokenManager: token,
			}
			tt.buildStubs(userRepo, todoRepo)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/todo-list/todo/%s", server.URL, tt.uri.ID)
			// body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			tt.setupAuth(request, tt.userID, token)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_doneTodo(t *testing.T) {

	tests := []struct {
		name          string
		uri           domain.TodoURI
		userID        string
		setupAuth     func(request *http.Request, id string, token auth.TokenManager)
		buildStubs    func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID)
		checkResponse func(recoder *httptest.ResponseRecorder, err1, err2 error)
	}{
		// TODO: Add test cases.
		{
			name: "status ok",
			uri: domain.TodoURI{
				ID: primitive.NewObjectID().Hex(),
			},
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID) {
				todoRepo.EXPECT().UpdateTodoDoneByID(gomock.Any(), todoID, userID).Return(domain.Todo{
					ID:       todoID,
					UserID:   userID,
					Title:    utils.RandomString(10),
					ActiveAt: utils.RandomString(10),
					Author:   utils.RandomString(10),
					Status:   domain.Done,
				}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err1, err2 error) {
				require.Equal(t, http.StatusOK, recoder.Code)
				require.NoError(t, err1)
				require.NoError(t, err2)
			},
		},
		{
			name: "status ok",
			uri: domain.TodoURI{
				ID: "primitive.NewObjectID().Hex()",
			},
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID) {
				todoRepo.EXPECT().UpdateTodoDoneByID(gomock.Any(), todoID, userID).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err1, err2 error) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
				require.Error(t, err1)
				require.NoError(t, err2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepo := repoMocks.NewMockTodo(ctrl)
			userRepo := repoMocks.NewMockUsers(ctrl)
			todoService := service.NewTodoService(todoRepo, userRepo)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			token, err := auth.NewManager("qwerty")
			require.NoError(t, err)
			handler := &Server{
				todoService:  todoService,
				tokenManager: token,
			}
			todoID, err1 := primitive.ObjectIDFromHex(tt.uri.ID)
			userID, err2 := primitive.ObjectIDFromHex(tt.userID)

			tt.buildStubs(userRepo, todoRepo, todoID, userID)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/todo-list/todo/%s/done", server.URL, tt.uri.ID)
			// body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPut, url, nil)
			tt.setupAuth(request, tt.userID, token)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder, err1, err2)
		})
	}
}

func TestServer_getTodos(t *testing.T) {

	tests := []struct {
		name          string
		userID        string
		status        string
		setupAuth     func(request *http.Request, id string, token auth.TokenManager)
		buildStubs    func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID, status string)
		checkResponse func(recoder *httptest.ResponseRecorder, err2 error)
	}{
		// TODO: Add test cases.
		{
			name:   "status ok",
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			status: domain.Active,
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID, status string) {
				todoRepo.EXPECT().GetTodoByStatus(gomock.Any(), status, userID).Return([]domain.Todo{
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					}}, nil)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err2 error) {
				require.Equal(t, http.StatusOK, recoder.Code)
				require.NoError(t, err2)
			},
		},
		{
			name:   "status ok",
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			status: domain.Done,
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID, status string) {
				todoRepo.EXPECT().GetTodoByStatus(gomock.Any(), status, userID).Return([]domain.Todo{
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					}}, nil)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err2 error) {
				require.Equal(t, http.StatusOK, recoder.Code)
				require.NoError(t, err2)
			},
		},
		{
			name:   "user id errors",
			userID: "primitive.NewObjectID().Hex()",
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			status: domain.Done,
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID, status string) {
				todoRepo.EXPECT().GetTodoByStatus(gomock.Any(), status, userID).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err2 error) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
				require.Error(t, err2)
			},
		},
		{
			name:   "user id errors",
			userID: primitive.NewObjectID().Hex(),
			setupAuth: func(request *http.Request, id string, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", id, time.Minute)
			},
			status: domain.Done,
			buildStubs: func(userRepo *repoMocks.MockUsers, todoRepo *repoMocks.MockTodo, todoID, userID primitive.ObjectID, status string) {
				todoRepo.EXPECT().GetTodoByStatus(gomock.Any(), status, userID).Times(1).Return([]domain.Todo{}, domain.ErrInternalServer)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder, err2 error) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
				require.NoError(t, err2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepo := repoMocks.NewMockTodo(ctrl)
			userRepo := repoMocks.NewMockUsers(ctrl)
			todoService := service.NewTodoService(todoRepo, userRepo)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			token, err := auth.NewManager("qwerty")
			require.NoError(t, err)
			handler := &Server{
				todoService:  todoService,
				tokenManager: token,
			}
			userID, err2 := primitive.ObjectIDFromHex(tt.userID)

			tt.buildStubs(userRepo, todoRepo, primitive.NewObjectID(), userID, tt.status)

			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()
			url := fmt.Sprintf("%s/api/v1/users/todo-list/todo?status=%s", server.URL, tt.status)
			// body, err := json.Marshal(tt.inp)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			tt.setupAuth(request, tt.userID, token)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder, err2)
		})
	}
}
