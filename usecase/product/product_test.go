package product

import (
	"errors"
	"product-service/models/product"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Insert(product *product.RegisterRequest) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetList(request *product.GetListRequest) (*[]product.Product, *[]product.GetAvailableStock, int, error) {
	args := m.Called(request)
	return args.Get(0).(*[]product.Product), args.Get(1).(*[]product.GetAvailableStock), args.Int(2), args.Error(3)
}

// Mock StockHttpRepository
type MockStockHttpRepository struct {
	mock.Mock
}

func (m *MockStockHttpRepository) GetAvailableStock(productList *[]product.GetAvailableStock) (map[int]int, error) {
	args := m.Called(productList)
	return args.Get(0).(map[int]int), args.Error(1)
}

func TestGetList_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockStockRepo := new(MockStockHttpRepository)
	usecase := NewProductUsecase(mockProductRepo, mockStockRepo)

	// Mock request
	req := &product.GetListRequest{
		Category: "Electronics",
		ShopId:   1,
		PerPage:  10,
		Page:     1,
	}

	// Mock products data
	mockProducts := &[]product.Product{
		{Id: 1, Name: "Laptop", Category: "Electronics", Price: 1000, ShopId: 1},
		{Id: 2, Name: "Mouse", Category: "Electronics", Price: 50, ShopId: 1},
	}

	// Mock ProductShopList
	mockProductShopList := &[]product.GetAvailableStock{
		{ProductId: 1, ShopId: 1},
		{ProductId: 2, ShopId: 1},
	}

	// Mock Stock Data
	mockStock := map[int]int{
		1: 5,  // 5 units in stock for ProductId 1
		2: 10, // 10 units in stock for ProductId 2
	}

	// Mock GetList
	mockProductRepo.On("GetList", req).Return(mockProducts, mockProductShopList, 2, nil)
	mockStockRepo.On("GetAvailableStock", mockProductShopList).Return(mockStock, nil)

	// Run Usecase
	products, total, err := usecase.GetList(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, *products, 2)
	assert.Equal(t, 5, (*products)[0].Stock)
	assert.Equal(t, 10, (*products)[1].Stock)

	// Verify mocks
	mockProductRepo.AssertExpectations(t)
	mockStockRepo.AssertExpectations(t)
}

func TestGetList_ProductRepoError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockStockRepo := new(MockStockHttpRepository)
	usecase := NewProductUsecase(mockProductRepo, mockStockRepo)

	req := &product.GetListRequest{Category: "Electronics", ShopId: 1, PerPage: 10, Page: 1}

	// Mock GetList error
	mockProductRepo.On("GetList", req).Return((*[]product.Product)(nil), (*[]product.GetAvailableStock)(nil), 0, errors.New("database error"))

	products, total, err := usecase.GetList(req)

	assert.Error(t, err)
	assert.Nil(t, products)
	assert.Equal(t, 0, total)

	mockProductRepo.AssertExpectations(t)
}

func TestGetList_StockRepoError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockStockRepo := new(MockStockHttpRepository)
	usecase := NewProductUsecase(mockProductRepo, mockStockRepo)

	req := &product.GetListRequest{Category: "Electronics", ShopId: 1, PerPage: 10, Page: 1}

	productsMock := []product.Product{
		{Id: 1, Name: "Laptop", Category: "Electronics", Price: 1000, ShopId: 1},
	}
	productShopListMock := []product.GetAvailableStock{
		{ProductId: 1, ShopId: 1},
	}

	mockProductRepo.On("GetList", req).Return(&productsMock, &productShopListMock, len(productsMock), nil)

	// Mock stock service error
	mockStockRepo.On("GetAvailableStock", &productShopListMock).Return(map[int]int(nil), errors.New("stock service error"))

	products, total, err := usecase.GetList(req)

	assert.Error(t, err)
	assert.Nil(t, products)
	assert.Equal(t, 0, total)

	mockProductRepo.AssertExpectations(t)
	mockStockRepo.AssertExpectations(t)
}
