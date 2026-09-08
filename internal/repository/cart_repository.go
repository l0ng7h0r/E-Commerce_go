package repository

import (
	"database/sql"
	"errors"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateCart(userID string) (*domain.Cart, error) {
	cart := &domain.Cart{UserID: userID}
	query := `SELECT id, created_at, updated_at FROM carts WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertQuery := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, created_at, updated_at`
			err = r.db.QueryRow(insertQuery, userID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Fetch items
	itemQuery := `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
		       p.name, p.price, COALESCE(p.image_url, '')
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1`
	rows, err := r.db.Query(itemQuery, cart.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item domain.CartItem
			var p domain.Product
			if err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt, &p.Name, &p.Price, &p.ImageURL); err == nil {
				p.ID = item.ProductID
				item.Product = &p
				cart.Items = append(cart.Items, item)
			}
		}
	}
	return cart, nil
}

func (r *CartRepository) AddItem(cartID, productID string, quantity int) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW()`
	_, err := r.db.Exec(query, cartID, productID, quantity)
	return err
}

func (r *CartRepository) UpdateItem(cartID, productID string, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE cart_id = $2 AND product_id = $3`
	_, err := r.db.Exec(query, quantity, cartID, productID)
	return err
}

func (r *CartRepository) RemoveItem(cartID, productID string) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	_, err := r.db.Exec(query, cartID, productID)
	return err
}

func (r *CartRepository) ClearCart(cartID string) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(query, cartID)
	return err
}
