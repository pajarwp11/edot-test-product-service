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

func (p *ProductRepository) GetList(request *product.GetListRequest) (*[]product.Product, *[]product.GetAvailableStock, int, error) {
	products := []product.Product{}
	getAvailableStock := []product.GetAvailableStock{}
	query := "SELECT id, name, category, price, shop_id FROM products WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM products WHERE 1=1"
	params := map[string]interface{}{}

	if request.Category != "" {
		query += " AND category = :category"
		countQuery += " AND category = :category"
		params["category"] = request.Category
	}

	if request.ShopId != 0 {
		query += " AND shop_id = :shop_id"
		countQuery += " AND shop_id = :shop_id"
		params["shop_id"] = request.ShopId
	}

	queryCount, args, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, nil, 0, err
	}
	var totalItems int
	queryCount = p.mysql.Rebind(queryCount)
	err = p.mysql.Get(&totalItems, queryCount, args...)
	if err != nil {
		return nil, nil, 0, err
	}

	params["limit"] = request.PerPage
	params["offset"] = request.PerPage * (request.Page - 1)
	query += " LIMIT :limit OFFSET :offset"

	query, args, err = sqlx.Named(query, params)
	if err != nil {
		return nil, nil, 0, err
	}
	query = p.mysql.Rebind(query)

	rows, err := p.mysql.Queryx(query, args...)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var prod product.Product
		if err := rows.StructScan(&prod); err != nil {
			return nil, nil, 0, err
		}
		products = append(products, prod)

		stock := product.GetAvailableStock{
			ProductId: prod.Id,
			ShopId:    prod.ShopId,
		}
		getAvailableStock = append(getAvailableStock, stock)
	}

	return &products, &getAvailableStock, totalItems, nil
}
