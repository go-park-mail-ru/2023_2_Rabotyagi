// Code generated by MockGen. DO NOT EDIT.
// Source: internal/product/usecases/basket_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/product/usecases/basket_service.go -destination=internal/product/mocks/basket_service.go --package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockIBasketStorage is a mock of IBasketStorage interface.
type MockIBasketStorage struct {
	ctrl     *gomock.Controller
	recorder *MockIBasketStorageMockRecorder
}

// MockIBasketStorageMockRecorder is the mock recorder for MockIBasketStorage.
type MockIBasketStorageMockRecorder struct {
	mock *MockIBasketStorage
}

// NewMockIBasketStorage creates a new mock instance.
func NewMockIBasketStorage(ctrl *gomock.Controller) *MockIBasketStorage {
	mock := &MockIBasketStorage{ctrl: ctrl}
	mock.recorder = &MockIBasketStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBasketStorage) EXPECT() *MockIBasketStorageMockRecorder {
	return m.recorder
}

// AddOrderInBasket mocks base method.
func (m *MockIBasketStorage) AddOrderInBasket(ctx context.Context, userID, productID uint64, count uint32) (*models.OrderInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrderInBasket", ctx, userID, productID, count)
	ret0, _ := ret[0].(*models.OrderInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrderInBasket indicates an expected call of AddOrderInBasket.
func (mr *MockIBasketStorageMockRecorder) AddOrderInBasket(ctx, userID, productID, count any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrderInBasket", reflect.TypeOf((*MockIBasketStorage)(nil).AddOrderInBasket), ctx, userID, productID, count)
}

// BuyFullBasket mocks base method.
func (m *MockIBasketStorage) BuyFullBasket(ctx context.Context, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyFullBasket", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyFullBasket indicates an expected call of BuyFullBasket.
func (mr *MockIBasketStorageMockRecorder) BuyFullBasket(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyFullBasket", reflect.TypeOf((*MockIBasketStorage)(nil).BuyFullBasket), ctx, userID)
}

// DeleteOrder mocks base method.
func (m *MockIBasketStorage) DeleteOrder(ctx context.Context, orderID, ownerID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", ctx, orderID, ownerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockIBasketStorageMockRecorder) DeleteOrder(ctx, orderID, ownerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockIBasketStorage)(nil).DeleteOrder), ctx, orderID, ownerID)
}

// GetOrdersInBasketByUserID mocks base method.
func (m *MockIBasketStorage) GetOrdersInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersInBasketByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.OrderInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersInBasketByUserID indicates an expected call of GetOrdersInBasketByUserID.
func (mr *MockIBasketStorageMockRecorder) GetOrdersInBasketByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersInBasketByUserID", reflect.TypeOf((*MockIBasketStorage)(nil).GetOrdersInBasketByUserID), ctx, userID)
}

// UpdateOrderCount mocks base method.
func (m *MockIBasketStorage) UpdateOrderCount(ctx context.Context, userID, orderID uint64, newCount uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderCount", ctx, userID, orderID, newCount)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderCount indicates an expected call of UpdateOrderCount.
func (mr *MockIBasketStorageMockRecorder) UpdateOrderCount(ctx, userID, orderID, newCount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderCount", reflect.TypeOf((*MockIBasketStorage)(nil).UpdateOrderCount), ctx, userID, orderID, newCount)
}

// UpdateOrderStatus mocks base method.
func (m *MockIBasketStorage) UpdateOrderStatus(ctx context.Context, userID, orderID uint64, newStatus uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, userID, orderID, newStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockIBasketStorageMockRecorder) UpdateOrderStatus(ctx, userID, orderID, newStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockIBasketStorage)(nil).UpdateOrderStatus), ctx, userID, orderID, newStatus)
}
