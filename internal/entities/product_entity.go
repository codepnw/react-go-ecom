package entities

import "time"

type Product struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Sold        int       `json:"sold"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductPayloadReq struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Price       float32   `json:"price" binding:"required"`
	Sold        int       `json:"sold"`
	Quantity    int       `json:"quantity" binding:"required"`
}
