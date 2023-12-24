// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/product/delivery/product_handler.go
//
// Generated by this command:
//
//	mockgen --source=./internal/product/delivery/product_handler.go --destination=internal/product/mocks/product_handler.go --package=mocks
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

// MockIProductService is a mock of IProductService interface.
type MockIProductService struct {
	ctrl     *gomock.Controller
	recorder *MockIProductServiceMockRecorder
}

// MockIProductServiceMockRecorder is the mock recorder for MockIProductService.
type MockIProductServiceMockRecorder struct {
	mock *MockIProductService
}

// NewMockIProductService creates a new mock instance.
func NewMockIProductService(ctrl *gomock.Controller) *MockIProductService {
	mock := &MockIProductService{ctrl: ctrl}
	mock.recorder = &MockIProductServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProductService) EXPECT() *MockIProductServiceMockRecorder {
	return m.recorder
}

// ActivateProduct mocks base method.
func (m *MockIProductService) ActivateProduct(ctx context.Context, productID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateProduct", ctx, productID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateProduct indicates an expected call of ActivateProduct.
func (mr *MockIProductServiceMockRecorder) ActivateProduct(ctx, productID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateProduct", reflect.TypeOf((*MockIProductService)(nil).ActivateProduct), ctx, productID, userID)
}

// AddComment mocks base method.
func (m *MockIProductService) AddComment(ctx context.Context, r io.Reader, userID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", ctx, r, userID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddComment indicates an expected call of AddComment.
func (mr *MockIProductServiceMockRecorder) AddComment(ctx, r, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockIProductService)(nil).AddComment), ctx, r, userID)
}

// AddOrder mocks base method.
func (m *MockIProductService) AddOrder(ctx context.Context, r io.Reader, userID uint64) (*models.OrderInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, r, userID)
	ret0, _ := ret[0].(*models.OrderInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockIProductServiceMockRecorder) AddOrder(ctx, r, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockIProductService)(nil).AddOrder), ctx, r, userID)
}

// AddPremium mocks base method.
func (m *MockIProductService) AddPremium(ctx context.Context, productID, userID, periodCode uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPremium", ctx, productID, userID, periodCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPremium indicates an expected call of AddPremium.
func (mr *MockIProductServiceMockRecorder) AddPremium(ctx, productID, userID, periodCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPremium", reflect.TypeOf((*MockIProductService)(nil).AddPremium), ctx, productID, userID, periodCode)
}

// AddProduct mocks base method.
func (m *MockIProductService) AddProduct(ctx context.Context, r io.Reader, userID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, r, userID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockIProductServiceMockRecorder) AddProduct(ctx, r, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockIProductService)(nil).AddProduct), ctx, r, userID)
}

// AddToFavourites mocks base method.
func (m *MockIProductService) AddToFavourites(ctx context.Context, userID uint64, r io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToFavourites", ctx, userID, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToFavourites indicates an expected call of AddToFavourites.
func (mr *MockIProductServiceMockRecorder) AddToFavourites(ctx, userID, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToFavourites", reflect.TypeOf((*MockIProductService)(nil).AddToFavourites), ctx, userID, r)
}

// BuyFullBasket mocks base method.
func (m *MockIProductService) BuyFullBasket(ctx context.Context, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyFullBasket", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyFullBasket indicates an expected call of BuyFullBasket.
func (mr *MockIProductServiceMockRecorder) BuyFullBasket(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyFullBasket", reflect.TypeOf((*MockIProductService)(nil).BuyFullBasket), ctx, userID)
}

// CloseProduct mocks base method.
func (m *MockIProductService) CloseProduct(ctx context.Context, productID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseProduct", ctx, productID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseProduct indicates an expected call of CloseProduct.
func (mr *MockIProductServiceMockRecorder) CloseProduct(ctx, productID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseProduct", reflect.TypeOf((*MockIProductService)(nil).CloseProduct), ctx, productID, userID)
}

// DeleteComment mocks base method.
func (m *MockIProductService) DeleteComment(ctx context.Context, commentID, senderID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", ctx, commentID, senderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockIProductServiceMockRecorder) DeleteComment(ctx, commentID, senderID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockIProductService)(nil).DeleteComment), ctx, commentID, senderID)
}

// DeleteFromFavourites mocks base method.
func (m *MockIProductService) DeleteFromFavourites(ctx context.Context, userID, productID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromFavourites", ctx, userID, productID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromFavourites indicates an expected call of DeleteFromFavourites.
func (mr *MockIProductServiceMockRecorder) DeleteFromFavourites(ctx, userID, productID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromFavourites", reflect.TypeOf((*MockIProductService)(nil).DeleteFromFavourites), ctx, userID, productID)
}

// DeleteOrder mocks base method.
func (m *MockIProductService) DeleteOrder(ctx context.Context, orderID, ownerID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", ctx, orderID, ownerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockIProductServiceMockRecorder) DeleteOrder(ctx, orderID, ownerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockIProductService)(nil).DeleteOrder), ctx, orderID, ownerID)
}

// DeleteProduct mocks base method.
func (m *MockIProductService) DeleteProduct(ctx context.Context, productID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", ctx, productID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockIProductServiceMockRecorder) DeleteProduct(ctx, productID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockIProductService)(nil).DeleteProduct), ctx, productID, userID)
}

// GetCommentList mocks base method.
func (m *MockIProductService) GetCommentList(ctx context.Context, offset, count, userID uint64) ([]*models.CommentInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentList", ctx, offset, count, userID)
	ret0, _ := ret[0].([]*models.CommentInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentList indicates an expected call of GetCommentList.
func (mr *MockIProductServiceMockRecorder) GetCommentList(ctx, offset, count, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentList", reflect.TypeOf((*MockIProductService)(nil).GetCommentList), ctx, offset, count, userID)
}

// GetOrdersByUserID mocks base method.
func (m *MockIProductService) GetOrdersByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.OrderInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUserID indicates an expected call of GetOrdersByUserID.
func (mr *MockIProductServiceMockRecorder) GetOrdersByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUserID", reflect.TypeOf((*MockIProductService)(nil).GetOrdersByUserID), ctx, userID)
}

// GetOrdersNotInBasketByUserID mocks base method.
func (m *MockIProductService) GetOrdersNotInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderNotInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersNotInBasketByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.OrderNotInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersNotInBasketByUserID indicates an expected call of GetOrdersNotInBasketByUserID.
func (mr *MockIProductServiceMockRecorder) GetOrdersNotInBasketByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersNotInBasketByUserID", reflect.TypeOf((*MockIProductService)(nil).GetOrdersNotInBasketByUserID), ctx, userID)
}

// GetOrdersSoldByUserID mocks base method.
func (m *MockIProductService) GetOrdersSoldByUserID(ctx context.Context, userID uint64) ([]*models.OrderNotInBasket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersSoldByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.OrderNotInBasket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersSoldByUserID indicates an expected call of GetOrdersSoldByUserID.
func (mr *MockIProductServiceMockRecorder) GetOrdersSoldByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersSoldByUserID", reflect.TypeOf((*MockIProductService)(nil).GetOrdersSoldByUserID), ctx, userID)
}

// GetProduct mocks base method.
func (m *MockIProductService) GetProduct(ctx context.Context, productID, userID uint64) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", ctx, productID, userID)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockIProductServiceMockRecorder) GetProduct(ctx, productID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockIProductService)(nil).GetProduct), ctx, productID, userID)
}

// GetProductsList mocks base method.
func (m *MockIProductService) GetProductsList(ctx context.Context, offset, count, userID uint64) ([]*models.ProductInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsList", ctx, offset, count, userID)
	ret0, _ := ret[0].([]*models.ProductInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsList indicates an expected call of GetProductsList.
func (mr *MockIProductServiceMockRecorder) GetProductsList(ctx, offset, count, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsList", reflect.TypeOf((*MockIProductService)(nil).GetProductsList), ctx, offset, count, userID)
}

// GetProductsOfSaler mocks base method.
func (m *MockIProductService) GetProductsOfSaler(ctx context.Context, offset, count, userID uint64, isMy bool) ([]*models.ProductInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsOfSaler", ctx, offset, count, userID, isMy)
	ret0, _ := ret[0].([]*models.ProductInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsOfSaler indicates an expected call of GetProductsOfSaler.
func (mr *MockIProductServiceMockRecorder) GetProductsOfSaler(ctx, offset, count, userID, isMy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsOfSaler", reflect.TypeOf((*MockIProductService)(nil).GetProductsOfSaler), ctx, offset, count, userID, isMy)
}

// GetSearchProductFeed mocks base method.
func (m *MockIProductService) GetSearchProductFeed(ctx context.Context, searchInput string, lastNumber, limit, userID uint64) ([]*models.ProductInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSearchProductFeed", ctx, searchInput, lastNumber, limit, userID)
	ret0, _ := ret[0].([]*models.ProductInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSearchProductFeed indicates an expected call of GetSearchProductFeed.
func (mr *MockIProductServiceMockRecorder) GetSearchProductFeed(ctx, searchInput, lastNumber, limit, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSearchProductFeed", reflect.TypeOf((*MockIProductService)(nil).GetSearchProductFeed), ctx, searchInput, lastNumber, limit, userID)
}

// GetUserFavourites mocks base method.
func (m *MockIProductService) GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFavourites", ctx, userID)
	ret0, _ := ret[0].([]*models.ProductInFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFavourites indicates an expected call of GetUserFavourites.
func (mr *MockIProductServiceMockRecorder) GetUserFavourites(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFavourites", reflect.TypeOf((*MockIProductService)(nil).GetUserFavourites), ctx, userID)
}

// SearchProduct mocks base method.
func (m *MockIProductService) SearchProduct(ctx context.Context, searchInput string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchProduct", ctx, searchInput)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchProduct indicates an expected call of SearchProduct.
func (mr *MockIProductServiceMockRecorder) SearchProduct(ctx, searchInput any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchProduct", reflect.TypeOf((*MockIProductService)(nil).SearchProduct), ctx, searchInput)
}

// UpdateComment mocks base method.
func (m *MockIProductService) UpdateComment(ctx context.Context, r io.Reader, userID, commentID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", ctx, r, userID, commentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateComment indicates an expected call of UpdateComment.
func (mr *MockIProductServiceMockRecorder) UpdateComment(ctx, r, userID, commentID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockIProductService)(nil).UpdateComment), ctx, r, userID, commentID)
}

// UpdateOrderCount mocks base method.
func (m *MockIProductService) UpdateOrderCount(ctx context.Context, r io.Reader, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderCount", ctx, r, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderCount indicates an expected call of UpdateOrderCount.
func (mr *MockIProductServiceMockRecorder) UpdateOrderCount(ctx, r, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderCount", reflect.TypeOf((*MockIProductService)(nil).UpdateOrderCount), ctx, r, userID)
}

// UpdateOrderStatus mocks base method.
func (m *MockIProductService) UpdateOrderStatus(ctx context.Context, r io.Reader, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, r, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockIProductServiceMockRecorder) UpdateOrderStatus(ctx, r, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockIProductService)(nil).UpdateOrderStatus), ctx, r, userID)
}

// UpdateProduct mocks base method.
func (m *MockIProductService) UpdateProduct(ctx context.Context, r io.Reader, isPartialUpdate bool, productID, userAuthID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", ctx, r, isPartialUpdate, productID, userAuthID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockIProductServiceMockRecorder) UpdateProduct(ctx, r, isPartialUpdate, productID, userAuthID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockIProductService)(nil).UpdateProduct), ctx, r, isPartialUpdate, productID, userAuthID)
}
