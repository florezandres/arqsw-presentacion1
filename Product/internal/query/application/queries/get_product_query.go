package queries

type GetProductQuery struct {
	ProductID string
}
type ListProductsQuery struct {
	Page  int
	Limit int
}
