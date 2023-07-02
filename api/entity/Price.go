package entity

import "time"

type Price struct {
	ID        string    `json:"id"`
	Product   string    `json:"product_id" validate:"required"`
	Currency  string    `json:"currency" validate:"required,currency"`
	Amount    *float64  `json:"price" validate:"required"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

func NewPrice() *Price {
	return &Price{}
}
