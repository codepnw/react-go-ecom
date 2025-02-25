package entities

import "time"

type Product struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float32    `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProductPayloadReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}
