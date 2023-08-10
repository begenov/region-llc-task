package service

import (
	"context"
	"testing"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	mocksRepo "github.com/begenov/region-llc-task/internal/repository/mocks"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTodoService_CreateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepo := mocksRepo.NewMockTodo(ctrl)
	userRepo := mocksRepo.NewMockUsers(ctrl)

	todoService := NewTodoService(todoRepo, userRepo)

	type args struct {
		ctx  context.Context
		todo domain.Todo
		user domain.User
	}
	tests := []struct {
		name          string
		args          args
		buildStubs    func(user domain.User, todo domain.Todo)
		checkResponse func(todo domain.Todo, err error)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(10),
					ActiveAt: time.Now().Add(time.Hour * 52).Format(domain.Format),
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(1).Return(domain.User{
					ID:       todo.UserID,
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
				}, nil)

				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(1).Return(domain.Todo{
					ID:       todo.ID,
					UserID:   todo.UserID,
					Title:    todo.Title,
					Author:   user.UserName,
					ActiveAt: todo.ActiveAt,
					Status:   todo.Status,
				}, nil)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, todo)
			},
		},
		{
			name: "invalid empty title",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    "",
					ActiveAt: "2023-08-09",
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(0)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInvalidTitle)
			},
		},
		{
			name: "header length exceeds 200 characters",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(256),
					ActiveAt: "2023-08-09",
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(0)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrHeaderLength)
			},
		},
		{
			name: "incorrect date format",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(50),
					ActiveAt: "invalid data format",
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(0)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrIncorrectDateFormat)
			},
		},
		{
			name: "not found",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(50),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, domain.ErrNotFound)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(0)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
		{
			name: "title already exists",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(50),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				var count int64 = 1
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(0)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrTitleAlreadyExists)
			},
		},
		{
			name: "not found user",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(50),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(1).Return(domain.User{}, domain.ErrNotFound)
				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(0)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
		{
			name: "internal server by create",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(50),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Status:   domain.Active,
				},
				user: domain.User{
					UserName: utils.RandomString(10),
					Email:    utils.RandomEmail(),
					Password: utils.RandomString(10),
				},
			},
			buildStubs: func(user domain.User, todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				userRepo.EXPECT().GetUserByID(gomock.Any(), todo.UserID).Times(1).Return(domain.User{
					ID:       todo.UserID,
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
				}, nil)

				todoRepo.EXPECT().Create(gomock.Any(), todo).Times(1).Return(domain.Todo{}, domain.ErrInternalServer)
			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.todo.Author = tt.args.user.UserName
			tt.buildStubs(tt.args.user, tt.args.todo)

			todo, err := todoService.CreateTodo(tt.args.ctx, tt.args.todo)

			tt.checkResponse(todo, err)
		})
	}
}

func TestTodoService_UpdateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepo := mocksRepo.NewMockTodo(ctrl)
	userRepo := mocksRepo.NewMockUsers(ctrl)

	todoService := NewTodoService(todoRepo, userRepo)

	type args struct {
		ctx  context.Context
		todo domain.Todo
	}
	tests := []struct {
		name          string
		args          args
		buildStubs    func(todo domain.Todo)
		checkResponse func(todo domain.Todo, err error)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(10),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(1).Return(nil)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(1).Return(domain.Todo{
					ID:       todo.ID,
					UserID:   todo.UserID,
					Title:    todo.Title,
					ActiveAt: todo.ActiveAt,
					Author:   todo.Author,
					Status:   todo.Status,
				}, nil)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, todo)
			},
		},
		{
			name: "inavalid title empty",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    "",
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInvalidTitle)
			},
		},
		{
			name: "header length",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(255),
					ActiveAt: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrHeaderLength)
			},
		},
		{
			name: "incorrect date format",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: "asf",
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(0)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrIncorrectDateFormat)
			},
		},
		{
			name: "check title",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, domain.ErrInternalServer)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
		{
			name: "check title",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 1
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(0)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrTitleAlreadyExists)
			},
		},
		{
			name: "todo not found",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(1).Return(domain.ErrNotFound)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
		{
			name: "todo internal server",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(1).Return(domain.ErrInternalServer)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(0)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
		{
			name: "get user by id",
			args: args{
				ctx: ctx,
				todo: domain.Todo{
					ID:       primitive.NewObjectID(),
					UserID:   primitive.NewObjectID(),
					Title:    utils.RandomString(25),
					ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
					Author:   utils.RandomString(10),
					Status:   domain.Active,
				},
			},
			buildStubs: func(todo domain.Todo) {
				var count int64 = 0
				todoRepo.EXPECT().GetCountByTitle(gomock.Any(), todo.Title, todo.UserID).Times(1).Return(count, nil)
				todoRepo.EXPECT().UpdateTodo(gomock.Any(), todo).Times(1).Return(nil)
				todoRepo.EXPECT().GetTodoByID(gomock.Any(), todo.ID).Times(1).Return(domain.Todo{}, domain.ErrInternalServer)

			},
			checkResponse: func(todo domain.Todo, err error) {
				require.Equal(t, err, domain.ErrInternalServer)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(tt.args.todo)

			todo, err := todoService.UpdateTodo(tt.args.ctx, tt.args.todo)

			tt.checkResponse(todo, err)
		})
	}
}

func TestTodoService_DeleteTodoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepo := mocksRepo.NewMockTodo(ctrl)
	userRepo := mocksRepo.NewMockUsers(ctrl)

	todoService := NewTodoService(todoRepo, userRepo)
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name          string
		args          args
		buildStubs    func(id primitive.ObjectID)
		checkResponse func(err error)
	}{
		// TODO: Add test cases.
		{
			name: "Ok",
			args: args{
				ctx: ctx,
				id:  primitive.NewObjectID().Hex(),
			},
			buildStubs: func(id primitive.ObjectID) {
				todoRepo.EXPECT().DeleteTodoByID(gomock.Any(), id).Return(nil)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "todo invalid id",
			args: args{
				ctx: ctx,
				id:  utils.RandomString(10),
			},
			buildStubs: func(id primitive.ObjectID) {
				todoRepo.EXPECT().DeleteTodoByID(gomock.Any(), id).Times(0)
			},
			checkResponse: func(err error) {
				require.Equal(t, err, domain.ErrTodoInvalidId)
			},
		},
		{
			name: "not found",
			args: args{
				ctx: ctx,
				id:  primitive.NewObjectID().Hex(),
			},
			buildStubs: func(id primitive.ObjectID) {
				todoRepo.EXPECT().DeleteTodoByID(gomock.Any(), id).Times(1).Return(domain.ErrNotFound)
			},
			checkResponse: func(err error) {
				require.Equal(t, err, domain.ErrNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := primitive.ObjectIDFromHex(tt.args.id)
			tt.buildStubs(id)

			err := todoService.DeleteTodoByID(tt.args.ctx, tt.args.id)

			tt.checkResponse(err)
		})
	}
}

func TestTodoService_UpdateTodoDoneByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepo := mocksRepo.NewMockTodo(ctrl)
	userRepo := mocksRepo.NewMockUsers(ctrl)

	todoService := NewTodoService(todoRepo, userRepo)

	type args struct {
		ctx    context.Context
		id     string
		userID primitive.ObjectID
	}
	tests := []struct {
		name          string
		args          args
		buildStubs    func(todoID, userID primitive.ObjectID)
		checkResponse func(todo domain.Todo, err1, err2 error)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				ctx:    ctx,
				id:     primitive.NewObjectID().Hex(),
				userID: primitive.NewObjectID(),
			},
			buildStubs: func(todoID, userID primitive.ObjectID) {
				todoRepo.EXPECT().UpdateTodoDoneByID(gomock.Any(), todoID, userID).Return(domain.Todo{
					ID:       todoID,
					UserID:   userID,
					Title:    utils.RandomString(10),
					ActiveAt: utils.RandomString(10),
					Author:   utils.RandomString(10),
					Status:   domain.Done,
				}, nil)
			},
			checkResponse: func(todo domain.Todo, err1, err2 error) {
				require.NoError(t, err1)
				require.NoError(t, err2)
				require.Equal(t, todo.Status, domain.Done)
			},
		},
		{
			name: "todo invalid id",
			args: args{
				ctx:    ctx,
				id:     "",
				userID: primitive.NewObjectID(),
			},
			buildStubs: func(todoID, userID primitive.ObjectID) {
				todoRepo.EXPECT().UpdateTodoDoneByID(gomock.Any(), todoID, userID).Times(0)
			},
			checkResponse: func(todo domain.Todo, err1, err2 error) {
				require.Equal(t, err1.Error(), "the provided hex string is not a valid ObjectID")
				require.Equal(t, err2, domain.ErrTodoInvalidId)
				// require.Equal(t, todo.Status, domain.Done)
			},
		},
		{
			name: "not found",
			args: args{
				ctx:    ctx,
				id:     primitive.NewObjectID().Hex(),
				userID: primitive.NewObjectID(),
			},
			buildStubs: func(todoID, userID primitive.ObjectID) {
				todoRepo.EXPECT().UpdateTodoDoneByID(gomock.Any(), todoID, userID).Times(1).Return(domain.Todo{}, domain.ErrNotFound)
			},
			checkResponse: func(todo domain.Todo, err1, err2 error) {
				require.NoError(t, err1)
				require.Equal(t, err2, domain.ErrNotFound)
				// require.Equal(t, todo.Status, domain.Done)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoID, err1 := primitive.ObjectIDFromHex(tt.args.id)
			tt.buildStubs(todoID, tt.args.userID)

			todo, err2 := todoService.UpdateTodoDoneByID(tt.args.ctx, tt.args.id, tt.args.userID)

			tt.checkResponse(todo, err1, err2)
		})
	}
}

func TestTodoService_GetTodosByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepo := mocksRepo.NewMockTodo(ctrl)
	userRepo := mocksRepo.NewMockUsers(ctrl)

	todoService := NewTodoService(todoRepo, userRepo)

	type args struct {
		ctx    context.Context
		status string
		userID primitive.ObjectID
	}
	tests := []struct {
		name          string
		args          args
		buildStubs    func(status string, userID primitive.ObjectID)
		checkResponse func(todos []domain.Todo, err error)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				ctx:    ctx,
				status: domain.Active,
				userID: primitive.NewObjectID(),
			},
			buildStubs: func(status string, userID primitive.ObjectID) {
				todoRepo.EXPECT().GetTodoByStatus(gomock.Any(), status, userID).Return([]domain.Todo{
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
					{
						ID:       primitive.NewObjectID(),
						UserID:   userID,
						Title:    utils.RandomString(10),
						ActiveAt: time.Now().Add(time.Hour * 24).Format(domain.Format),
						Author:   utils.RandomString(10),
						Status:   status,
					},
				}, nil)
			},
			checkResponse: func(todos []domain.Todo, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(tt.args.status, tt.args.userID)

			todos, err := todoService.GetTodosByStatus(tt.args.ctx, tt.args.status, tt.args.userID)

			tt.checkResponse(todos, err)
		})
	}
}
