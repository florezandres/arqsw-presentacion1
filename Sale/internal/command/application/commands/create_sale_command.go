package commands

type CreateSaleCommand struct {
	ProductID  string
	Quantity   int
	TotalPrice float64
	SaleDate   string
}
