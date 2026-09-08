package domain

import "time"

type OrderItem struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   *Product `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	Price     float64  `json:"price"`
}

type Order struct {
	ID             string      `json:"id"`
	UserID         string      `json:"user_id"`
	TotalAmount    float64     `json:"total_amount"`
	Status         string      `json:"status"`
	LogisticBranch string      `json:"logistic_branch,omitempty"`
	OrderItems     []OrderItem `json:"order_items"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type CreateOrderReq struct {
	LogisticBranch string `json:"logistic_branch"`
}

type UpdateOrderStatusReq struct {
	Status string `json:"status" binding:"required"`
}
