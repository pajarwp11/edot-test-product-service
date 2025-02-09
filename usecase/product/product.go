package product

import (
	"product-service/models/product"
	"strconv"
)

type ProductRepository interface {
	Insert(product *product.RegisterRequest) error
	GetList(request *product.GetListRequest) (*[]product.Product, *[]product.GetAvailableStock, int, error)
}

type StockHttpRepository interface {
	GetAvailableStock(productList *[]product.GetAvailableStock) (map[string]interface{}, error)
}

type ProductUsecase struct {
	productRepo         ProductRepository
	stockHttpRepository StockHttpRepository
}

func NewProductUsecase(productRepo ProductRepository, stockHttpRepository StockHttpRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo:         productRepo,
		stockHttpRepository: stockHttpRepository,
	}
}

func (p *ProductUsecase) Register(productRegister *product.RegisterRequest) error {
	// TO DO: check if shop id belong to user
	return p.productRepo.Insert(productRegister)
}

func (p *ProductUsecase) GetList(request *product.GetListRequest) (*[]product.Product, int, error) {
	products, productShopList, total, err := p.productRepo.GetList(request)
	if err != nil {
		return nil, 0, err
	}
	if len(*products) > 0 {
		productStockMap, err := p.stockHttpRepository.GetAvailableStock(productShopList)
		if err != nil {
			return nil, 0, err
		}

		for i, product := range *products {
			productId := strconv.Itoa(product.Id)
			stock := productStockMap[productId]
			if stock != nil {
				stockProduct := stock.(float64)
				(*products)[i].Stock = int(stockProduct)
			}
		}
	}

	return products, total, nil
}
