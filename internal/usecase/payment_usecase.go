package usecase

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
	"github.com/l0ng7h0r/ecommerce/pkg/phajay"
)

type PaymentUsecase struct {
	paymentRepo  *repository.PaymentRepository
	orderRepo    *repository.OrderRepository
	phajayClient *phajay.Client
}

func NewPaymentUsecase(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository, phajayClient *phajay.Client) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		phajayClient: phajayClient,
	}
}

// ────────────────────────────────────────────────────────────────────────────
// CreatePayment – legacy redirect-link flow (kept for backward compat.)
// ────────────────────────────────────────────────────────────────────────────

func (u *PaymentUsecase) CreatePayment(userID string, req *domain.CreatePaymentReq) (*domain.PaymentResponse, error) {
	order, err := u.orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized order payment")
	}

	if order.Status == "paid" || order.Status == "confirmed" {
		return nil, errors.New("order is already paid")
	}

	// Return existing pending payment link if available
	existingPayment, _ := u.paymentRepo.GetPaymentByOrderID(order.ID)
	if existingPayment != nil && existingPayment.PaymentURL != "" && existingPayment.Status == "pending" {
		return &domain.PaymentResponse{
			PaymentID:  existingPayment.ID,
			OrderID:    existingPayment.OrderID,
			Amount:     existingPayment.Amount,
			Status:     existingPayment.Status,
			PaymentURL: existingPayment.PaymentURL,
		}, nil
	}

	description := fmt.Sprintf("Order %s", order.ID[:8])

	phajayResp, err := u.phajayClient.CreatePaymentLink(order.TotalAmount, description, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Phajay payment link: %w", err)
	}

	payment := &domain.Payment{
		OrderID:    order.ID,
		Amount:     order.TotalAmount,
		Status:     "pending",
		Method:     "phajay",
		PaymentURL: phajayResp.PaymentURL,
	}

	createdPayment, err := u.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	return &domain.PaymentResponse{
		PaymentID:  createdPayment.ID,
		OrderID:    createdPayment.OrderID,
		Amount:     createdPayment.Amount,
		Status:     createdPayment.Status,
		PaymentURL: phajayResp.PaymentURL,
	}, nil
}

// ────────────────────────────────────────────────────────────────────────────
// GenerateQR – new in-app QR Code flow
// ────────────────────────────────────────────────────────────────────────────

func (u *PaymentUsecase) GenerateQR(userID string, req *domain.GenerateQRReq) (*domain.GenerateQRResponse, error) {
	order, err := u.orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized order payment")
	}

	if order.Status == "paid" || order.Status == "confirmed" {
		return nil, errors.New("order is already paid")
	}

	// If there's already a pending payment with a QR code, reuse it
	existing, _ := u.paymentRepo.GetPaymentByOrderID(order.ID)
	if existing != nil && existing.Status == "pending" && existing.QRCode != "" {
		return &domain.GenerateQRResponse{
			PaymentID:     existing.ID,
			OrderID:       existing.OrderID,
			Amount:        existing.Amount,
			Status:        existing.Status,
			TransactionID: existing.TransactionID,
			QRCode:        existing.QRCode,
		}, nil
	}

	bank := req.Bank
	if bank == "" {
		bank = "BCEL"
	}

	description := fmt.Sprintf("Order %s", order.ID[:8])

	qrResp, err := u.phajayClient.GenerateQR(bank, order.TotalAmount, description, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Phajay QR: %w", err)
	}

	payment := &domain.Payment{
		OrderID:       order.ID,
		Amount:        order.TotalAmount,
		Status:        "pending",
		Method:        "phajay_qr",
		TransactionID: qrResp.TransactionID,
		QRCode:        qrResp.QRCode,
	}

	createdPayment, err := u.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	return &domain.GenerateQRResponse{
		PaymentID:     createdPayment.ID,
		OrderID:       createdPayment.OrderID,
		Amount:        createdPayment.Amount,
		Status:        createdPayment.Status,
		TransactionID: qrResp.TransactionID,
		QRCode:        qrResp.QRCode,
		DeepLink:      qrResp.Link,
	}, nil
}

// ────────────────────────────────────────────────────────────────────────────
// CheckPaymentStatus – poll Phajay and update DB if now paid
// ────────────────────────────────────────────────────────────────────────────

func (u *PaymentUsecase) CheckPaymentStatus(orderID, userID string, roles []string) (*domain.PaymentStatusResponse, error) {
	payment, err := u.paymentRepo.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	for _, r := range roles {
		if r == "admin" {
			isAdmin = true
			break
		}
	}
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("unauthorized to view this payment")
	}

	// If already paid, just return current status
	if payment.Status == "paid" || payment.Status == "completed" {
		return &domain.PaymentStatusResponse{
			PaymentID:     payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Status:        payment.Status,
			TransactionID: payment.TransactionID,
			QRCode:        payment.QRCode,
		}, nil
	}

	// Fetch live status from Phajay if we have a transaction ID
	if payment.TransactionID != "" {
		phajayStatus, err := u.phajayClient.CheckTransactionStatus(payment.TransactionID)
		if err != nil {
			log.Printf("[CheckPaymentStatus] error checking Phajay status for tx %s: %v", payment.TransactionID, err)
		} else {
			log.Printf("[CheckPaymentStatus] tx %s isPaid=%v status=%s", payment.TransactionID, phajayStatus.IsPaymentPaid(), phajayStatus.GetStatus())
			if phajayStatus.IsPaymentPaid() {
				_ = u.paymentRepo.UpdatePaymentStatus(payment.ID, "paid", payment.TransactionID)
				_ = u.orderRepo.UpdateOrderStatus(orderID, "paid")
				payment.Status = "paid"
			} else {
				normalized := strings.ToUpper(phajayStatus.GetStatus())
				if normalized == "FAILED" || normalized == "EXPIRED" || normalized == "CANCELLED" {
					_ = u.paymentRepo.UpdatePaymentStatus(payment.ID, "failed", payment.TransactionID)
					_ = u.orderRepo.UpdateOrderStatus(orderID, "cancelled")
					_ = u.orderRepo.RestoreStockForOrder(orderID)
					payment.Status = "failed"
				}
			}
		}
	}

	return &domain.PaymentStatusResponse{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		QRCode:        payment.QRCode,
	}, nil
}

// ────────────────────────────────────────────────────────────────────────────
// HandleWebhook – called by Phajay when payment completes / fails
// ────────────────────────────────────────────────────────────────────────────

func (u *PaymentUsecase) HandleWebhook(payload *phajay.WebhookPayload) error {
	log.Printf("[Phajay Webhook Received] %+v", payload)

	var payment *domain.Payment

	if payload.OrderNo != "" {
		var err error
		payment, err = u.paymentRepo.GetPaymentByOrderID(payload.OrderNo)
		if err != nil {
			log.Printf("[Webhook] lookup by OrderNo failed: %v", err)
		}
	}

	if payment == nil && payload.TransactionID != "" {
		var err error
		payment, err = u.paymentRepo.GetPaymentByTransactionID(payload.TransactionID)
		if err != nil {
			log.Printf("[Webhook] lookup by TransactionID failed: %v", err)
		}
	}

	if payment == nil {
		return fmt.Errorf("payment not found for order: %s / tx: %s", payload.OrderNo, payload.TransactionID)
	}

	if payload.IsPaymentPaid() {
		if err := u.paymentRepo.UpdatePaymentStatus(payment.ID, "paid", payload.TransactionID); err != nil {
			return err
		}
		return u.orderRepo.UpdateOrderStatus(payment.OrderID, "paid")
	}

	normalized := strings.ToUpper(payload.Status)
	if normalized == "FAILED" || normalized == "CANCELLED" || normalized == "EXPIRED" {
		if err := u.paymentRepo.UpdatePaymentStatus(payment.ID, "failed", payload.TransactionID); err != nil {
			return err
		}
		_ = u.orderRepo.UpdateOrderStatus(payment.OrderID, "cancelled")
		return u.orderRepo.RestoreStockForOrder(payment.OrderID)
	}

	return nil
}

// ────────────────────────────────────────────────────────────────────────────
// GetPaymentByOrderID – for authenticated status view
// ────────────────────────────────────────────────────────────────────────────

func (u *PaymentUsecase) GetPaymentByOrderID(orderID, currentUserID string, roles []string) (*domain.Payment, error) {
	payment, err := u.paymentRepo.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	for _, r := range roles {
		if r == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin && order.UserID != currentUserID {
		return nil, errors.New("unauthorized to view this payment")
	}

	return payment, nil
}
