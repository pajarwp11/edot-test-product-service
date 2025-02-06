package product

type Product struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Category string `db:"category"`
	Price    int    `db:"price"`
	ShopId   int    `db:"shop_id"`
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
