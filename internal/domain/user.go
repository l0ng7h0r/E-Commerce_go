package domain

import "time"

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	ShopName  string    `json:"shop_name,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SellerRegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	ShopName string `json:"shop_name" binding:"required"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	ShopName string `json:"shop_name,omitempty"`
	Role     string `json:"role"` // Optional, default to "user"
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EditUserReq struct {
	ID       string `json:"id" binding:"required"`
	Email    string `json:"email" binding:"omitempty"`
	Name     string `json:"name" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
}

type EditSellerReq struct {
	ID       string `json:"id" binding:"required"`
	Email    string `json:"email" binding:"omitempty"`
	Name     string `json:"name" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	ShopName string `json:"shop_name" binding:"omitempty"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	UserID       string   `json:"user_id"`
	Email        string   `json:"email"`
	Name         string   `json:"name,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	ShopName     string   `json:"shop_name,omitempty"`
	Roles        []string `json:"roles"`
}
