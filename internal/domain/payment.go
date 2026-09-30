package domain

import "time"

type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Method        string    `json:"method"`
	PaymentURL    string    `json:"payment_url,omitempty"`
	QRCode        string    `json:"qr_code,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePaymentReq struct {
	OrderID string `json:"order_id" binding:"required"`
}

type GenerateQRReq struct {
	OrderID string `json:"order_id" binding:"required"`
	Bank    string `json:"bank,omitempty"` // "BCEL","JDB","LDB","IB","STB","MMONEY" – defaults to BCEL
}

type GenerateQRResponse struct {
	PaymentID     string  `json:"payment_id"`
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	QRCode        string  `json:"qr_code"`
	DeepLink      string  `json:"deep_link,omitempty"`
}

type PaymentResponse struct {
	PaymentID  string  `json:"payment_id"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	PaymentURL string  `json:"payment_url"`
}

type PaymentStatusResponse struct {
	PaymentID     string  `json:"payment_id"`
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id,omitempty"`
	QRCode        string  `json:"qr_code,omitempty"`
}
