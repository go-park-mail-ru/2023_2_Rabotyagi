// Code generated by MockGen. DO NOT EDIT.
// Source: internal/category/delivery/category_handler.go
//
// Generated by this command:
//
//	mockgen -source=internal/category/delivery/category_handler.go -destination=internal/category/mocks/service.go
//
// Package mock_delivery is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockICategoryService is a mock of ICategoryService interface.
type MockICategoryService struct {
	ctrl     *gomock.Controller
	recorder *MockICategoryServiceMockRecorder
}

// MockICategoryServiceMockRecorder is the mock recorder for MockICategoryService.
type MockICategoryServiceMockRecorder struct {
	mock *MockICategoryService
}

// NewMockICategoryService creates a new mock instance.
func NewMockICategoryService(ctrl *gomock.Controller) *MockICategoryService {
	mock := &MockICategoryService{ctrl: ctrl}
	mock.recorder = &MockICategoryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICategoryService) EXPECT() *MockICategoryServiceMockRecorder {
	return m.recorder
}

// GetFullCategories mocks base method.
func (m *MockICategoryService) GetFullCategories(ctx context.Context) ([]*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFullCategories", ctx)
	ret0, _ := ret[0].([]*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFullCategories indicates an expected call of GetFullCategories.
func (mr *MockICategoryServiceMockRecorder) GetFullCategories(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullCategories", reflect.TypeOf((*MockICategoryService)(nil).GetFullCategories), ctx)
}

// SearchCategory mocks base method.
func (m *MockICategoryService) SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCategory", ctx, searchInput)
	ret0, _ := ret[0].([]*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCategory indicates an expected call of SearchCategory.
func (mr *MockICategoryServiceMockRecorder) SearchCategory(ctx, searchInput any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCategory", reflect.TypeOf((*MockICategoryService)(nil).SearchCategory), ctx, searchInput)
}
