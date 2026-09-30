package phajay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

// ────────────────────────────────────────────────────────────────────────────
// Payment Link (old redirect-based flow)
// ────────────────────────────────────────────────────────────────────────────

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

// ────────────────────────────────────────────────────────────────────────────
// Direct QR Payment (in-app QR code – no redirect)
// ────────────────────────────────────────────────────────────────────────────

type GenerateQRRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	OrderNo     string  `json:"orderNo"`
	Tag1        string  `json:"tag1,omitempty"`
}

// GenerateQRResponse is the common structure returned by all
// bank-specific QR endpoints (BCEL, JDB, LDB, IB, STB, M-Money …)
type GenerateQRResponse struct {
	TransactionID string `json:"transactionId"`
	// QRCode is the raw EMVCo QR payload string the frontend renders
	QRCode string `json:"qrCode"`
	// Link is a deep-link / mobile URI for direct app launch
	Link string `json:"link"`
}

// bankQRPath maps a bank code to the Phajay API path suffix.
var bankQRPath = map[string]string{
	"BCEL":    "/payment/generate-bcel-qr",
	"JDB":     "/payment/generate-jdb-qr",
	"LDB":     "/payment/generate-ldb-qr",
	"IB":      "/payment/generate-ib-qr",
	"STB":     "/payment/generate-stb-qr",
	"MMONEY":  "/payment/generate-m-money-qr",
}

// GenerateQR calls the bank-specific Phajay QR endpoint and returns
// a raw QR payload that the frontend can render with qrcode.react.
// bank should be one of the keys in bankQRPath (e.g. "BCEL").
// Falls back to BCEL if the bank is unknown.
func (c *Client) GenerateQR(bank string, amount float64, description, orderNo string) (*GenerateQRResponse, error) {
	path, ok := bankQRPath[bank]
	if !ok {
		path = bankQRPath["BCEL"]
	}

	reqBody := GenerateQRRequest{
		Amount:      amount,
		Description: description,
		OrderNo:     orderNo,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QR request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create QR request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SecretKey, "")
	req.Header.Set("secretKey", c.SecretKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call phajay QR API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read QR response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("phajay QR API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result GenerateQRResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse QR response: %w", err)
	}

	return &result, nil
}

// ────────────────────────────────────────────────────────────────────────────
// Check Transaction Status (poll payment result from Phajay)
// ────────────────────────────────────────────────────────────────────────────

type TransactionData struct {
	TransactionID string  `json:"transactionId"`
	OrderNo       string  `json:"orderNo"`
	LinkCode      string  `json:"linkCode"`
	Status        string  `json:"status"` // e.g. "WAITING", "PAYMENT_COMPLETED", "EXPIRED", "REFUNDED"
	IsPaid        bool    `json:"isPaid"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	PaymentMethod string  `json:"paymentMethod"`
	PaymentTime   string  `json:"paymentTime"`
	CreatedAt     string  `json:"createdAt"`
}

type TransactionStatusResponse struct {
	Message string          `json:"message"`
	Data    TransactionData `json:"data"`
	// Direct flat fields fallback if ever returned
	TransactionID string  `json:"transactionId"`
	Status        string  `json:"status"`
	IsPaid        bool    `json:"isPaid"`
}

func (r *TransactionStatusResponse) GetStatus() string {
	if r.Data.Status != "" {
		return r.Data.Status
	}
	return r.Status
}

func (r *TransactionStatusResponse) IsPaymentPaid() bool {
	if r.Data.IsPaid || r.IsPaid {
		return true
	}
	st := strings.ToUpper(r.GetStatus())
	return st == "PAYMENT_COMPLETED" || st == "SUCCESS" || st == "COMPLETED" || st == "PAID" || st == "SUCCESSFULLY"
}

func (c *Client) CheckTransactionStatus(transactionID string) (*TransactionStatusResponse, error) {
	url := fmt.Sprintf("%s/payment/check-transaction/status/%s", c.BaseURL, transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	encodedKey := base64.StdEncoding.EncodeToString([]byte(c.SecretKey))
	req.Header.Set("Authorization", "Basic "+encodedKey)
	req.Header.Set("secretKey", c.SecretKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call phajay status API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read status response: %w", err)
	}

	log.Printf("[Phajay CheckTransactionStatus] txID: %s, code: %d, body: %s", transactionID, resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("phajay status API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result TransactionStatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	return &result, nil
}

// ────────────────────────────────────────────────────────────────────────────
// Webhook Payload (sent by Phajay to our /webhooks/phajay endpoint)
// ────────────────────────────────────────────────────────────────────────────

type WebhookPayload struct {
	OrderNo       string  `json:"orderNo"`
	TransactionID string  `json:"transactionId"`
	BillNumber    string  `json:"billNumber"`
	ExReferenceNo string  `json:"exReferenceNo"`
	Amount        float64 `json:"txnAmount"`
	AmountAlt     float64 `json:"amount"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	PaymentMethod string  `json:"paymentMethod"`
}

func (w *WebhookPayload) IsPaymentPaid() bool {
	st := strings.ToUpper(w.Status)
	msg := strings.ToUpper(w.Message)
	return st == "PAYMENT_COMPLETED" || st == "SUCCESS" || st == "COMPLETED" || st == "PAID" || msg == "SUCCESS"
}
