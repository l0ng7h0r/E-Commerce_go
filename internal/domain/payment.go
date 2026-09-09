package domain

import "time"

type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Method        string    `json:"method"`
	PaymentURL    string    `json:"payment_url,omitempty"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePaymentReq struct {
	OrderID string `json:"order_id" binding:"required"`
}

type PaymentResponse struct {
	PaymentID  string  `json:"payment_id"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	PaymentURL string  `json:"payment_url"`
}
