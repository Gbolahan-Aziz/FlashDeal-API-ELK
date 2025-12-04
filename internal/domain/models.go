package domain

import "time"

type Deal struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
	Active bool    `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type NewDeal struct {
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
	Active bool    `json:"active"`
}

type Order struct {
	ID      string    `json:"order_id"`
	DealID  string    `json:"deal_id"`
	Qty     int       `json:"qty"`
	Status  string    `json:"status"`
	Created time.Time `json:"created_at"`
}

type NewOrder struct {
	DealID string `json:"deal_id"`
	Qty    int    `json:"qty"`
	UserID string `json:"user_id,omitempty"`
}
