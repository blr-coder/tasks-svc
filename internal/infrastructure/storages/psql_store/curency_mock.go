// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/storages/psql_store/curency.go

// Package psql_store is a generated GoMock package.
package psql_store

import (
	context "context"
	reflect "reflect"

	models "github.com/blr-coder/tasks-svc/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockICurrencyStorage is a mock of ICurrencyStorage interface.
type MockICurrencyStorage struct {
	ctrl     *gomock.Controller
	recorder *MockICurrencyStorageMockRecorder
}

// MockICurrencyStorageMockRecorder is the mock recorder for MockICurrencyStorage.
type MockICurrencyStorageMockRecorder struct {
	mock *MockICurrencyStorage
}

// NewMockICurrencyStorage creates a new mock instance.
func NewMockICurrencyStorage(ctrl *gomock.Controller) *MockICurrencyStorage {
	mock := &MockICurrencyStorage{ctrl: ctrl}
	mock.recorder = &MockICurrencyStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICurrencyStorage) EXPECT() *MockICurrencyStorageMockRecorder {
	return m.recorder
}

// GetRateByEUR mocks base method.
func (m *MockICurrencyStorage) GetRateByEUR(ctx context.Context, currency models.Currency) (*models.CurrencyRate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRateByEUR", ctx, currency)
	ret0, _ := ret[0].(*models.CurrencyRate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRateByEUR indicates an expected call of GetRateByEUR.
func (mr *MockICurrencyStorageMockRecorder) GetRateByEUR(ctx, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRateByEUR", reflect.TypeOf((*MockICurrencyStorage)(nil).GetRateByEUR), ctx, currency)
}

// ListCurrencyTickers mocks base method.
func (m *MockICurrencyStorage) ListCurrencyTickers(ctx context.Context) (models.CurrencyList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCurrencyTickers", ctx)
	ret0, _ := ret[0].(models.CurrencyList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCurrencyTickers indicates an expected call of ListCurrencyTickers.
func (mr *MockICurrencyStorageMockRecorder) ListCurrencyTickers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCurrencyTickers", reflect.TypeOf((*MockICurrencyStorage)(nil).ListCurrencyTickers), ctx)
}

// SetCurrencyRates mocks base method.
func (m *MockICurrencyStorage) SetCurrencyRates(ctx context.Context, rates []models.CurrencyRate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrencyRates", ctx, rates)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrencyRates indicates an expected call of SetCurrencyRates.
func (mr *MockICurrencyStorageMockRecorder) SetCurrencyRates(ctx, rates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrencyRates", reflect.TypeOf((*MockICurrencyStorage)(nil).SetCurrencyRates), ctx, rates)
}
