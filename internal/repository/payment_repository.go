package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(payment *domain.Payment) (*domain.Payment, error) {
	query := `
		INSERT INTO payments (order_id, amount, status, transaction_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	err := r.db.QueryRow(query, payment.OrderID, payment.Amount, payment.Status, payment.TransactionID).
		Scan(&payment.ID, &payment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return payment, nil
}

func (r *PaymentRepository) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, amount, status, COALESCE(transaction_id, ''), created_at FROM payments WHERE order_id = $1`
	p := &domain.Payment{}
	err := r.db.QueryRow(query, orderID).Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.TransactionID, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return p, nil
}
