// Code generated by MockGen. DO NOT EDIT.
// Source: internal/product/usecases/premium_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/product/usecases/premium_service.go -destination=internal/product/mocks/premium_service.go --package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockIPremiumStorage is a mock of IPremiumStorage interface.
type MockIPremiumStorage struct {
	ctrl     *gomock.Controller
	recorder *MockIPremiumStorageMockRecorder
}

// MockIPremiumStorageMockRecorder is the mock recorder for MockIPremiumStorage.
type MockIPremiumStorageMockRecorder struct {
	mock *MockIPremiumStorage
}

// NewMockIPremiumStorage creates a new mock instance.
func NewMockIPremiumStorage(ctrl *gomock.Controller) *MockIPremiumStorage {
	mock := &MockIPremiumStorage{ctrl: ctrl}
	mock.recorder = &MockIPremiumStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPremiumStorage) EXPECT() *MockIPremiumStorageMockRecorder {
	return m.recorder
}

// AddPremium mocks base method.
func (m *MockIPremiumStorage) AddPremium(ctx context.Context, productID, userID uint64, premiumBegin, premiumExpire time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPremium", ctx, productID, userID, premiumBegin, premiumExpire)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPremium indicates an expected call of AddPremium.
func (mr *MockIPremiumStorageMockRecorder) AddPremium(ctx, productID, userID, premiumBegin, premiumExpire any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPremium", reflect.TypeOf((*MockIPremiumStorage)(nil).AddPremium), ctx, productID, userID, premiumBegin, premiumExpire)
}
