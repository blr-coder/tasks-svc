// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/storages/transaction/transaction_manager.go

// Package transaction is a generated GoMock package.
package transaction

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockITransaction is a mock of ITransaction interface.
type MockITransaction struct {
	ctrl     *gomock.Controller
	recorder *MockITransactionMockRecorder
}

// MockITransactionMockRecorder is the mock recorder for MockITransaction.
type MockITransactionMockRecorder struct {
	mock *MockITransaction
}

// NewMockITransaction creates a new mock instance.
func NewMockITransaction(ctrl *gomock.Controller) *MockITransaction {
	mock := &MockITransaction{ctrl: ctrl}
	mock.recorder = &MockITransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITransaction) EXPECT() *MockITransactionMockRecorder {
	return m.recorder
}

// Finish mocks base method.
func (m *MockITransaction) Finish(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finish", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Finish indicates an expected call of Finish.
func (mr *MockITransactionMockRecorder) Finish(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finish", reflect.TypeOf((*MockITransaction)(nil).Finish), arg0)
}

// GetTx mocks base method.
func (m *MockITransaction) GetTx() *sqlx.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx")
	ret0, _ := ret[0].(*sqlx.Tx)
	return ret0
}

// GetTx indicates an expected call of GetTx.
func (mr *MockITransactionMockRecorder) GetTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockITransaction)(nil).GetTx))
}

// Rollback mocks base method.
func (m *MockITransaction) Rollback(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockITransactionMockRecorder) Rollback(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockITransaction)(nil).Rollback), arg0)
}

// MockIDBTransactionManager is a mock of IDBTransactionManager interface.
type MockIDBTransactionManager struct {
	ctrl     *gomock.Controller
	recorder *MockIDBTransactionManagerMockRecorder
}

// MockIDBTransactionManagerMockRecorder is the mock recorder for MockIDBTransactionManager.
type MockIDBTransactionManagerMockRecorder struct {
	mock *MockIDBTransactionManager
}

// NewMockIDBTransactionManager creates a new mock instance.
func NewMockIDBTransactionManager(ctrl *gomock.Controller) *MockIDBTransactionManager {
	mock := &MockIDBTransactionManager{ctrl: ctrl}
	mock.recorder = &MockIDBTransactionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDBTransactionManager) EXPECT() *MockIDBTransactionManagerMockRecorder {
	return m.recorder
}

// StartTx mocks base method.
func (m *MockIDBTransactionManager) StartTx(arg0 context.Context) (ITransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTx", arg0)
	ret0, _ := ret[0].(ITransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartTx indicates an expected call of StartTx.
func (mr *MockIDBTransactionManagerMockRecorder) StartTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTx", reflect.TypeOf((*MockIDBTransactionManager)(nil).StartTx), arg0)
}
