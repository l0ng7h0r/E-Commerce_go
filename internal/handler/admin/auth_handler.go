package admin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type AdminAuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAdminAuthHandler(authUsecase *usecase.AuthUsecase) *AdminAuthHandler {
	return &AdminAuthHandler{authUsecase: authUsecase}
}

// Login godoc
// @Summary Admin Login
// @Description Login to Admin Management Portal
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginReq true "Login Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /admin/login [post]
func (h *AdminAuthHandler) Login(c fiber.Ctx) error {
	var req domain.LoginReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.authUsecase.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if admin
	hasRole := false
	for _, r := range res.Roles {
		if r == "admin" {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account does not have admin permissions"})
	}

	return c.JSON(res)
}

// EditUser godoc
// @Summary Edit User
// @Description Edit user profile (name, email, phone) by ID
// @Tags Admin User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body domain.EditUserReq true "Edit User Request"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [patch]
func (h *AdminAuthHandler) EditUser(c fiber.Ctx) error {
	id := c.Params("id")
	var req domain.EditUserReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if id != "" {
		req.ID = id
	}
	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	res, err := h.authUsecase.EditUser(req.ID, req.Email, req.Name, req.Phone)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// EditSeller godoc
// @Summary Edit Seller
// @Description Edit seller profile (name, email, phone, shop_name) by ID
// @Tags Admin User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Seller User ID"
// @Param request body domain.EditSellerReq true "Edit Seller Request"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/sellers/{id} [patch]
func (h *AdminAuthHandler) EditSeller(c fiber.Ctx) error {
	id := c.Params("id")
	var req domain.EditSellerReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if id != "" {
		req.ID = id
	}
	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Seller ID is required"})
	}

	res, err := h.authUsecase.EditSeller(req.ID, req.Email, req.Name, req.Phone, req.ShopName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// Refresh godoc
// @Summary Admin Refresh Token
// @Description Refresh admin access token
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /admin/refresh [post]
func (h *AdminAuthHandler) Refresh(c fiber.Ctx) error {
	var req domain.RefreshTokenReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.authUsecase.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// Logout godoc
// @Summary Admin Logout
// @Description Revoke admin refresh token
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} map[string]string
// @Router /admin/logout [post]
func (h *AdminAuthHandler) Logout(c fiber.Ctx) error {
	var req domain.RefreshTokenReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_ = h.authUsecase.Logout(req.RefreshToken)
	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}
