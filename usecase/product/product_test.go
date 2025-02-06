package product

import (
	"errors"
	"product-service/models/product"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Insert(product *product.RegisterRequest) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetList(request *product.GetListRequest) (*[]product.Product, int, error) {
	args := m.Called(request)
	if args.Get(0) != nil {
		return args.Get(0).(*[]product.Product), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

// Test Register Success
func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	mockProduct := &product.RegisterRequest{
		Name:     "Test Product",
		Category: "Electronics",
		Price:    1000,
	}

	mockRepo.On("Insert", mockProduct).Return(nil)

	err := usecase.Register(mockProduct)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test Register Failure
func TestRegister_Failure(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	mockProduct := &product.RegisterRequest{
		Name:     "Test Product",
		Category: "Electronics",
		Price:    1000,
	}

	mockRepo.On("Insert", mockProduct).Return(errors.New("database error"))

	err := usecase.Register(mockProduct)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test GetList Success
func TestGetList_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	request := &product.GetListRequest{Category: "Electronics", Page: 1, PerPage: 10}
	mockProducts := []product.Product{
		{Id: 1, Name: "Product A", Category: "Electronics", Price: 1000},
		{Id: 2, Name: "Product B", Category: "Electronics", Price: 2000},
	}

	mockRepo.On("GetList", request).Return(&mockProducts, 2, nil)

	products, total, err := usecase.GetList(request)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(*products))
	mockRepo.AssertExpectations(t)
}

// Test GetList Failure
func TestGetList_Failure(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	request := &product.GetListRequest{Category: "Electronics", Page: 1, PerPage: 10}

	mockRepo.On("GetList", request).Return(nil, 0, errors.New("database error"))

	products, total, err := usecase.GetList(request)
	assert.Error(t, err)
	assert.Nil(t, products)
	assert.Equal(t, 0, total)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}
