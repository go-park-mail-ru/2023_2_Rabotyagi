// Code generated by MockGen. DO NOT EDIT.
// Source: delivery/favourite_handler.go
//
// Generated by this command:
//
//	mockgen -source=delivery/favourite_handler.go -destination=mocks/favourite_handler.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockIFavouriteService is a mock of IFavouriteService interface.
type MockIFavouriteService struct {
	ctrl     *gomock.Controller
	recorder *MockIFavouriteServiceMockRecorder
}

// MockIFavouriteServiceMockRecorder is the mock recorder for MockIFavouriteService.
type MockIFavouriteServiceMockRecorder struct {
	mock *MockIFavouriteService
}

// NewMockIFavouriteService creates a new mock instance.
func NewMockIFavouriteService(ctrl *gomock.Controller) *MockIFavouriteService {
	mock := &MockIFavouriteService{ctrl: ctrl}
	mock.recorder = &MockIFavouriteServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFavouriteService) EXPECT() *MockIFavouriteServiceMockRecorder {
	return m.recorder
}

// AddToFavourites mocks base method.
func (m *MockIFavouriteService) AddToFavourites(ctx context.Context, userID uint64, r io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToFavourites", ctx, userID, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToFavourites indicates an expected call of AddToFavourites.
func (mr *MockIFavouriteServiceMockRecorder) AddToFavourites(ctx, userID, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToFavourites", reflect.TypeOf((*MockIFavouriteService)(nil).AddToFavourites), ctx, userID, r)
}

// DeleteFromFavourites mocks base method.
func (m *MockIFavouriteService) DeleteFromFavourites(ctx context.Context, userID, productID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromFavourites", ctx, userID, productID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromFavourites indicates an expected call of DeleteFromFavourites.
func (mr *MockIFavouriteServiceMockRecorder) DeleteFromFavourites(ctx, userID, productID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromFavourites", reflect.TypeOf((*MockIFavouriteService)(nil).DeleteFromFavourites), ctx, userID, productID)
}

// GetUserFavourites mocks base method.
func (m *MockIFavouriteService) GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFavourites", ctx, userID)
	ret0, _ := ret[0].([]*models.ProductInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFavourites indicates an expected call of GetUserFavourites.
func (mr *MockIFavouriteServiceMockRecorder) GetUserFavourites(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFavourites", reflect.TypeOf((*MockIFavouriteService)(nil).GetUserFavourites), ctx, userID)
}
