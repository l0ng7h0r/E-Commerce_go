package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, hashedPassword string) (*domain.User, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at, updated_at`
	user := &domain.User{}
	err := r.db.QueryRow(query, email, hashedPassword).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) AssignRole(userID, roleName string) error {
	var roleID string
	err := r.db.QueryRow(`SELECT id FROM roles WHERE name = $1`, roleName).Scan(&roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Insert role if not exists
			err = r.db.QueryRow(`INSERT INTO roles (name) VALUES ($1) RETURNING id`, roleName).Scan(&roleID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err = r.db.Exec(query, userID, roleID)
	return err
}

func (r *UserRepository) GetUserRoles(userID string) ([]string, error) {
	query := `
		SELECT r.name FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	roles, err := r.GetUserRoles(user.ID)
	if err == nil {
		user.Roles = roles
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	roles, err := r.GetUserRoles(user.ID)
	if err == nil {
		user.Roles = roles
	}
	return user, nil
}

func (r *UserRepository) SaveRefreshToken(userID, token string) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, NOW() + INTERVAL '7 days')`
	_, err := r.db.Exec(query, userID, token)
	return err
}

func (r *UserRepository) GetRefreshToken(token string) (string, error) {
	var userID string
	query := `SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid or expired refresh token")
		}
		return "", err
	}
	return userID, nil
}

func (r *UserRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *UserRepository) GetAllUsers() ([]*domain.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		roles, _ := r.GetUserRoles(u.ID)
		u.Roles = roles
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
