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
