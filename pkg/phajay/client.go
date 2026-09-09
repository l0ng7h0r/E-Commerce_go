package phajay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	SecretKey string
	BaseURL   string
	client    *http.Client
}

func NewClient(secretKey string) *Client {
	return &Client{
		SecretKey: secretKey,
		BaseURL:   "https://payment-gateway.phajay.co/v1/api",
		client:    &http.Client{},
	}
}

type PaymentLinkRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	OrderNo     string  `json:"orderNo"`
}

type PaymentLinkResponse struct {
	Success    bool   `json:"success"`
	PaymentURL string `json:"redirectURL"`
	Message    string `json:"message"`
}

type WebhookPayload struct {
	OrderNo       string  `json:"orderNo"`
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"` // success, failed, cancelled
}

func (c *Client) CreatePaymentLink(amount float64, description, orderNo string) (*PaymentLinkResponse, error) {
	reqBody := PaymentLinkRequest{
		Amount:      amount,
		Description: description,
		OrderNo:     orderNo,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/link/payment-link", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SecretKey, "")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call phajay API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("phajay API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result PaymentLinkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
