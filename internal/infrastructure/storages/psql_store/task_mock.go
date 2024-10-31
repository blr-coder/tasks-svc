// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/storages/psql_store/task.go

// Package psql_store is a generated GoMock package.
package psql_store

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	models "github.com/blr-coder/tasks-svc/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockITaskStorage is a mock of ITaskStorage interface.
type MockITaskStorage struct {
	ctrl     *gomock.Controller
	recorder *MockITaskStorageMockRecorder
}

// MockITaskStorageMockRecorder is the mock recorder for MockITaskStorage.
type MockITaskStorageMockRecorder struct {
	mock *MockITaskStorage
}

// NewMockITaskStorage creates a new mock instance.
func NewMockITaskStorage(ctrl *gomock.Controller) *MockITaskStorage {
	mock := &MockITaskStorage{ctrl: ctrl}
	mock.recorder = &MockITaskStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITaskStorage) EXPECT() *MockITaskStorageMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockITaskStorage) Count(ctx context.Context, filter *models.TasksFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filter)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockITaskStorageMockRecorder) Count(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockITaskStorage)(nil).Count), ctx, filter)
}

// Create mocks base method.
func (m *MockITaskStorage) Create(ctx context.Context, createTask *models.CreateTask) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, createTask)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockITaskStorageMockRecorder) Create(ctx, createTask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockITaskStorage)(nil).Create), ctx, createTask)
}

// Delete mocks base method.
func (m *MockITaskStorage) Delete(ctx context.Context, taskID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockITaskStorageMockRecorder) Delete(ctx, taskID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockITaskStorage)(nil).Delete), ctx, taskID)
}

// Get mocks base method.
func (m *MockITaskStorage) Get(ctx context.Context, taskID int64) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, taskID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockITaskStorageMockRecorder) Get(ctx, taskID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockITaskStorage)(nil).Get), ctx, taskID)
}

// List mocks base method.
func (m *MockITaskStorage) List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockITaskStorageMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockITaskStorage)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockITaskStorage) Update(ctx context.Context, input *models.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, input)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockITaskStorageMockRecorder) Update(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockITaskStorage)(nil).Update), ctx, input)
}

// WithTransaction mocks base method.
func (m *MockITaskStorage) WithTransaction(tx *sqlx.Tx) ITaskStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", tx)
	ret0, _ := ret[0].(ITaskStorage)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockITaskStorageMockRecorder) WithTransaction(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockITaskStorage)(nil).WithTransaction), tx)
}

// MockIStorageExecutor is a mock of IStorageExecutor interface.
type MockIStorageExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockIStorageExecutorMockRecorder
}

// MockIStorageExecutorMockRecorder is the mock recorder for MockIStorageExecutor.
type MockIStorageExecutorMockRecorder struct {
	mock *MockIStorageExecutor
}

// NewMockIStorageExecutor creates a new mock instance.
func NewMockIStorageExecutor(ctrl *gomock.Controller) *MockIStorageExecutor {
	mock := &MockIStorageExecutor{ctrl: ctrl}
	mock.recorder = &MockIStorageExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStorageExecutor) EXPECT() *MockIStorageExecutorMockRecorder {
	return m.recorder
}

// ExecContext mocks base method.
func (m *MockIStorageExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockIStorageExecutorMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockIStorageExecutor)(nil).ExecContext), varargs...)
}

// GetContext mocks base method.
func (m *MockIStorageExecutor) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockIStorageExecutorMockRecorder) GetContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockIStorageExecutor)(nil).GetContext), varargs...)
}

// QueryContext mocks base method.
func (m *MockIStorageExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockIStorageExecutorMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockIStorageExecutor)(nil).QueryContext), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockIStorageExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockIStorageExecutorMockRecorder) QueryRowContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockIStorageExecutor)(nil).QueryRowContext), varargs...)
}

// QueryRowxContext mocks base method.
func (m *MockIStorageExecutor) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowxContext", varargs...)
	ret0, _ := ret[0].(*sqlx.Row)
	return ret0
}

// QueryRowxContext indicates an expected call of QueryRowxContext.
func (mr *MockIStorageExecutorMockRecorder) QueryRowxContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowxContext", reflect.TypeOf((*MockIStorageExecutor)(nil).QueryRowxContext), varargs...)
}

// SelectContext mocks base method.
func (m *MockIStorageExecutor) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectContext indicates an expected call of SelectContext.
func (mr *MockIStorageExecutorMockRecorder) SelectContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContext", reflect.TypeOf((*MockIStorageExecutor)(nil).SelectContext), varargs...)
}