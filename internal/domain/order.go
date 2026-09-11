package domain

import "time"

type OrderItem struct {
	ID        string       `json:"id"`
	OrderID   string       `json:"order_id"`
	ProductID string       `json:"product_id"`
	Product   *CartProduct `json:"product,omitempty"`
	Quantity  int          `json:"quantity"`
	Price     float64      `json:"price"`
}

type Order struct {
	ID             string      `json:"id"`
	UserID         string      `json:"user_id"`
	TotalAmount    float64     `json:"total_amount"`
	Status         string      `json:"status"`
	LogisticCompany string      `json:"logistic_company,omitempty"`
	LogisticBranch string      `json:"logistic_branch,omitempty"`
	District       string      `json:"district,omitempty"`
	OrderItems     []OrderItem `json:"order_items"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type CreateOrderReq struct {
	LogisticCompany string `json:"logistic_company"`
	LogisticBranch  string `json:"logistic_branch"`
	District    string `json:"district"`
}

type UpdateOrderStatusReq struct {
	Status string `json:"status" binding:"required"`
}
