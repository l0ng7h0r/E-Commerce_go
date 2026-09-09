package usecase

import (
	"errors"
	"fmt"

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

	// Check if existing pending payment link exists
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

	// Call Phajay API to create payment link
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

func (u *PaymentUsecase) HandleWebhook(payload *phajay.WebhookPayload) error {
	payment, err := u.paymentRepo.GetPaymentByOrderID(payload.OrderNo)
	if err != nil {
		return fmt.Errorf("payment not found for order: %s", payload.OrderNo)
	}

	switch payload.Status {
	case "success":
		if err := u.paymentRepo.UpdatePaymentStatus(payment.ID, "paid", payload.TransactionID); err != nil {
			return err
		}
		return u.orderRepo.UpdateOrderStatus(payload.OrderNo, "paid")
	case "failed", "cancelled":
		if err := u.paymentRepo.UpdatePaymentStatus(payment.ID, "failed", payload.TransactionID); err != nil {
			return err
		}
		_ = u.orderRepo.UpdateOrderStatus(payload.OrderNo, "cancelled")
		// Restore stock back for cancelled/failed order
		return u.orderRepo.RestoreStockForOrder(payload.OrderNo)
	default:
		return fmt.Errorf("unknown webhook status: %s", payload.Status)
	}
}

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
