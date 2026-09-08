package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Category methods
func (r *ProductRepository) CreateCategory(name string) (*domain.Category, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at`
	cat := &domain.Category{}
	err := r.db.QueryRow(query, name).Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return cat, nil
}

func (r *ProductRepository) GetAllCategories() ([]*domain.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories ORDER BY name ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Category
	for rows.Next() {
		cat := &domain.Category{}
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, cat)
	}
	return list, nil
}

func (r *ProductRepository) GetCategoryByID(id string) (*domain.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	cat := &domain.Category{}
	err := r.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return cat, nil
}

// Product methods
func (r *ProductRepository) CreateProduct(p *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (seller_id, category_id, name, description, price, stock, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, p.SellerID, p.CategoryID, p.Name, p.Description, p.Price, p.Stock, p.ImageURL).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) GetAllProducts() ([]*domain.Product, error) {
	query := `SELECT id, seller_id, COALESCE(category_id, ''), name, description, price, stock, COALESCE(image_url, ''), created_at, updated_at FROM products ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ProductRepository) GetProductByID(id string) (*domain.Product, error) {
	query := `SELECT id, seller_id, COALESCE(category_id, ''), name, description, price, stock, COALESCE(image_url, ''), created_at, updated_at FROM products WHERE id = $1`
	p := &domain.Product{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) GetProductsBySellerID(sellerID string) ([]*domain.Product, error) {
	query := `SELECT id, seller_id, COALESCE(category_id, ''), name, description, price, stock, COALESCE(image_url, ''), created_at, updated_at FROM products WHERE seller_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ProductRepository) UpdateProduct(p *domain.Product) error {
	query := `
		UPDATE products
		SET category_id = $1, name = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = NOW()
		WHERE id = $7 AND seller_id = $8`
	res, err := r.db.Exec(query, p.CategoryID, p.Name, p.Description, p.Price, p.Stock, p.ImageURL, p.ID, p.SellerID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("product not found or unauthorized")
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id, sellerID string) error {
	query := `DELETE FROM products WHERE id = $1 AND seller_id = $2`
	res, err := r.db.Exec(query, id, sellerID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("product not found or unauthorized")
	}
	return nil
}
