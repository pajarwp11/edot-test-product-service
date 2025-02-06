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
