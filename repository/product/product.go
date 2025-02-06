package product

import (
	"product-service/models/product"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	mysql *sqlx.DB
}

func NewProductRepository(mysql *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		mysql: mysql,
	}
}

func (p *ProductRepository) Insert(product *product.RegisterRequest) error {
	_, err := p.mysql.Exec("INSERT INTO products (name,category,price,shop_id) VALUES (?,?,?,?)", product.Name, product.Category, product.Price, product.ShopId)
	return err
}

func (p *ProductRepository) GetList(request *product.GetListRequest) (*[]product.Product, error) {
	products := []product.Product{}
	query := "SELECT id, name, category, price, shop_id FROM products WHERE 1=1"
	params := map[string]interface{}{}

	if request.Category != "" {
		query += " AND category = :category"
		params["category"] = request.Category
	}

	if request.ShopId != 0 {
		query += " AND shop_id = :shop_id"
		params["shop_id"] = request.ShopId
	}

	limit := request.PerPage
	offset := limit * (request.Page - 1)

	query += " LIMTT = :limit OFFSET = :offset"
	params["limit"] = limit
	params["offset"] = offset

	stmt, err := p.mysql.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	err = stmt.Select(&products, params)
	if err != nil {
		return nil, err
	}

	return &products, err
}
