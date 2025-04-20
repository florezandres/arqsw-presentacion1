package queries

type GetSaleQuery struct {
	SaleID string
}

type ListSalesQuery struct {
	Page  int
	Limit int
}
