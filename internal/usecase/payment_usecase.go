package usecase

import (
	"errors"

	"github.com/google/uuid"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
)

type PaymentUsecase struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
}

func NewPaymentUsecase(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (u *PaymentUsecase) CreatePayment(userID string, req *domain.CreatePaymentReq) (*domain.Payment, error) {
	order, err := u.orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized order payment")
	}

	if order.Status == "paid" {
		return nil, errors.New("order is already paid")
	}

	payment := &domain.Payment{
		OrderID:       order.ID,
		Amount:        order.TotalAmount,
		Status:        "completed",
		TransactionID: "TXN-" + uuid.New().String()[:8],
	}

	createdPayment, err := u.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	_ = u.orderRepo.UpdateOrderStatus(order.ID, "paid")

	return createdPayment, nil
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
