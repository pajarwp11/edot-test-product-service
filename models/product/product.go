package product

type Product struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	Price    int    `db:"price" json:"price"`
	ShopId   int    `db:"shop_id" json:"shop_id"`
	Stock    int    `json:"stock"`
}

type GetListRequest struct {
	Category string
	ShopId   int
	Page     int
	PerPage  int
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Category string `json:"category" validate:"required"`
	Price    int    `json:"price" validate:"required"`
	ShopId   int    `json:"shop_id" validate:"required"`
}

type GetAvailableStock struct {
	ProductId int `json:"product_id"`
	ShopId    int `json:"shop_id"`
}
