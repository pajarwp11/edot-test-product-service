package product

import (
	"product-service/models/product"
)

type ProductRepository interface {
	Insert(product *product.RegisterRequest) error
	GetList(request *product.GetListRequest) (*[]product.Product, *[]product.GetAvailableStock, int, error)
}

type ProductUsecase struct {
	productRepo ProductRepository
}

func NewProductUsecase(productRepo ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (p *ProductUsecase) Register(productRegister *product.RegisterRequest) error {
	return p.productRepo.Insert(productRegister)
}

func (p *ProductUsecase) GetList(request *product.GetListRequest) (*[]product.Product, int, error) {
	products, _, total, err := p.productRepo.GetList(request)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
