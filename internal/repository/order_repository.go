package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order, items []domain.OrderItem) (*domain.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO orders (user_id, total_amount, status, logistic_branch)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err = tx.QueryRow(query, order.UserID, order.TotalAmount, order.Status, order.LogisticBranch).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`
	for i := range items {
		items[i].OrderID = order.ID
		err = tx.QueryRow(itemQuery, items[i].OrderID, items[i].ProductID, items[i].Quantity, items[i].Price).Scan(&items[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	order.OrderItems = items
	return order, nil
}

func (r *OrderRepository) GetOrderByID(id string) (*domain.Order, error) {
	query := `SELECT id, user_id, total_amount, status, COALESCE(logistic_branch, ''), created_at, updated_at FROM orders WHERE id = $1`
	order := &domain.Order{}
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.LogisticBranch, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
		       COALESCE(p.name, ''), COALESCE(p.image_url, '')
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1`
	rows, err := r.db.Query(itemsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item domain.OrderItem
			var p domain.Product
			if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &p.Name, &p.ImageURL); err == nil {
				p.ID = item.ProductID
				item.Product = &p
				order.OrderItems = append(order.OrderItems, item)
			}
		}
	}
	return order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID string) ([]*domain.Order, error) {
	query := `SELECT id, user_id, total_amount, status, COALESCE(logistic_branch, ''), created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.LogisticBranch, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	for _, o := range orders {
		itemsQuery := `
			SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
			       COALESCE(p.name, ''), COALESCE(p.image_url, '')
			FROM order_items oi
			LEFT JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = $1`
		itemRows, err := r.db.Query(itemsQuery, o.ID)
		if err == nil {
			for itemRows.Next() {
				var item domain.OrderItem
				var p domain.Product
				if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &p.Name, &p.ImageURL); err == nil {
					p.ID = item.ProductID
					item.Product = &p
					o.OrderItems = append(o.OrderItems, item)
				}
			}
			itemRows.Close()
		}
	}
	return orders, nil
}

func (r *OrderRepository) GetAllOrders() ([]*domain.Order, error) {
	query := `SELECT id, user_id, total_amount, status, COALESCE(logistic_branch, ''), created_at, updated_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.LogisticBranch, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("order not found")
	}
	return nil
}
