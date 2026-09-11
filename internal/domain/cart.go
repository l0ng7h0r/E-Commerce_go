package domain

import "time"

type CartItem struct {
	ID        string       `json:"id"`
	CartID    string       `json:"cart_id"`
	ProductID string       `json:"product_id"`
	Product   *CartProduct `json:"product,omitempty"`
	Quantity  int          `json:"quantity"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartProduct struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
}

type AddCartItemReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type UpdateCartItemReq struct {
	Quantity int `json:"quantity" binding:"required"`
}
