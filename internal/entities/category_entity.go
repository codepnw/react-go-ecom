package entities

import "time"

type Category struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryReq struct {
	Title string `form:"title" json:"title" binding:"required"`
}
