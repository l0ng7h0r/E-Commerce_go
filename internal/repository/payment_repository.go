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
		INSERT INTO payments (order_id, amount, status, transaction_id, method, payment_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	err := r.db.QueryRow(query, payment.OrderID, payment.Amount, payment.Status, payment.TransactionID, payment.Method, payment.PaymentURL).
		Scan(&payment.ID, &payment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return payment, nil
}

func (r *PaymentRepository) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, amount, status, COALESCE(transaction_id, ''), COALESCE(method, 'phajay'), COALESCE(payment_url, ''), created_at FROM payments WHERE order_id = $1`
	p := &domain.Payment{}
	err := r.db.QueryRow(query, orderID).Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.TransactionID, &p.Method, &p.PaymentURL, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return p, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(paymentID, status, transactionID string) error {
	query := `UPDATE payments SET status = $1, transaction_id = $2 WHERE id = $3`
	res, err := r.db.Exec(query, status, transactionID, paymentID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("payment not found")
	}
	return nil
}
