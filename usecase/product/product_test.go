package product

import (
	"errors"
	"product-service/models/product"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for ProductRepository
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

// Test case for GetList
func TestGetList(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	mockProducts := []product.Product{
		{Id: 1, Name: "Product A", Category: "Category 1", Price: 100, ShopId: 1},
		{Id: 2, Name: "Product B", Category: "Category 2", Price: 200, ShopId: 2},
	}
	mockStock := []product.GetAvailableStock{
		{ProductId: 1, ShopId: 10},
		{ProductId: 2, ShopId: 20},
	}
	total := 2

	request := &product.GetListRequest{
		Category: "Category 1",
		ShopId:   1,
		Page:     1,
		PerPage:  10,
	}

	// Define expected return values
	mockRepo.On("GetList", request).Return(&mockProducts, &mockStock, total, nil)

	// Call the usecase method
	products, totalItems, err := usecase.GetList(request)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, total, totalItems)
	assert.Equal(t, mockProducts, *products)

	// Ensure that the mock expectations were met
	mockRepo.AssertExpectations(t)
}

// Test case when repository returns an error
func TestGetList_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	request := &product.GetListRequest{
		Category: "Invalid Category",
		ShopId:   1,
		Page:     1,
		PerPage:  10,
	}

	mockRepo.On("GetList", request).Return(&[]product.Product{}, &[]product.GetAvailableStock{}, 0, errors.New("database error"))

	products, totalItems, err := usecase.GetList(request)

	assert.Error(t, err)
	assert.Equal(t, 0, totalItems)
	assert.Nil(t, products)

	mockRepo.AssertExpectations(t)
}
