package models

type ProductRead struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int32
}
