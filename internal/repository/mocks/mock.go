// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/begenov/region-llc-task/internal/domain"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsers) Create(ctx context.Context, user domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsersMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsers)(nil).Create), ctx, user)
}

// GetByRefreshToken mocks base method.
func (m *MockUsers) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRefreshToken indicates an expected call of GetByRefreshToken.
func (mr *MockUsersMockRecorder) GetByRefreshToken(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRefreshToken", reflect.TypeOf((*MockUsers)(nil).GetByRefreshToken), ctx, refreshToken)
}

// GetUserByEmail mocks base method.
func (m *MockUsers) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUsersMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUsers)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockUsers) GetUserByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUsersMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUsers)(nil).GetUserByID), ctx, id)
}

// SetSession mocks base method.
func (m *MockUsers) SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", ctx, userID, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockUsersMockRecorder) SetSession(ctx, userID, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockUsers)(nil).SetSession), ctx, userID, session)
}

// MockTodo is a mock of Todo interface.
type MockTodo struct {
	ctrl     *gomock.Controller
	recorder *MockTodoMockRecorder
}

// MockTodoMockRecorder is the mock recorder for MockTodo.
type MockTodoMockRecorder struct {
	mock *MockTodo
}

// NewMockTodo creates a new mock instance.
func NewMockTodo(ctrl *gomock.Controller) *MockTodo {
	mock := &MockTodo{ctrl: ctrl}
	mock.recorder = &MockTodoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodo) EXPECT() *MockTodoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodo) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, todo)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoMockRecorder) Create(ctx, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodo)(nil).Create), ctx, todo)
}

// DeleteTodoByID mocks base method.
func (m *MockTodo) DeleteTodoByID(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTodoByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTodoByID indicates an expected call of DeleteTodoByID.
func (mr *MockTodoMockRecorder) DeleteTodoByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTodoByID", reflect.TypeOf((*MockTodo)(nil).DeleteTodoByID), ctx, id)
}

// GetCountByTitle mocks base method.
func (m *MockTodo) GetCountByTitle(ctx context.Context, title string, id primitive.ObjectID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountByTitle", ctx, title, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountByTitle indicates an expected call of GetCountByTitle.
func (mr *MockTodoMockRecorder) GetCountByTitle(ctx, title, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountByTitle", reflect.TypeOf((*MockTodo)(nil).GetCountByTitle), ctx, title, id)
}

// GetTodoByID mocks base method.
func (m *MockTodo) GetTodoByID(ctx context.Context, id primitive.ObjectID) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodoByID", ctx, id)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodoByID indicates an expected call of GetTodoByID.
func (mr *MockTodoMockRecorder) GetTodoByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodoByID", reflect.TypeOf((*MockTodo)(nil).GetTodoByID), ctx, id)
}

// GetTodoByStatus mocks base method.
func (m *MockTodo) GetTodoByStatus(ctx context.Context, status string, userID primitive.ObjectID) ([]domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodoByStatus", ctx, status, userID)
	ret0, _ := ret[0].([]domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodoByStatus indicates an expected call of GetTodoByStatus.
func (mr *MockTodoMockRecorder) GetTodoByStatus(ctx, status, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodoByStatus", reflect.TypeOf((*MockTodo)(nil).GetTodoByStatus), ctx, status, userID)
}

// UpdateTodo mocks base method.
func (m *MockTodo) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodo", ctx, todo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTodo indicates an expected call of UpdateTodo.
func (mr *MockTodoMockRecorder) UpdateTodo(ctx, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodo", reflect.TypeOf((*MockTodo)(nil).UpdateTodo), ctx, todo)
}

// UpdateTodoDoneByID mocks base method.
func (m *MockTodo) UpdateTodoDoneByID(ctx context.Context, id, userID primitive.ObjectID) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodoDoneByID", ctx, id, userID)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTodoDoneByID indicates an expected call of UpdateTodoDoneByID.
func (mr *MockTodoMockRecorder) UpdateTodoDoneByID(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodoDoneByID", reflect.TypeOf((*MockTodo)(nil).UpdateTodoDoneByID), ctx, id, userID)
}

// UpdateTodoID mocks base method.
func (m *MockTodo) UpdateTodoID(ctx context.Context, todo domain.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodoID", ctx, todo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTodoID indicates an expected call of UpdateTodoID.
func (mr *MockTodoMockRecorder) UpdateTodoID(ctx, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodoID", reflect.TypeOf((*MockTodo)(nil).UpdateTodoID), ctx, todo)
}

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRedis) Delete(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRedisMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRedis)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockRedis) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedis)(nil).Get), key)
}

// Set mocks base method.
func (m *MockRedis) Set(key, value string, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisMockRecorder) Set(key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedis)(nil).Set), key, value, expiration)
}
