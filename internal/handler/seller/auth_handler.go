package seller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type SellerAuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewSellerAuthHandler(authUsecase *usecase.AuthUsecase) *SellerAuthHandler {
	return &SellerAuthHandler{authUsecase: authUsecase}
}

// Login godoc
// @Summary Seller Login
// @Description Login to Seller Portal
// @Tags Seller Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginReq true "Login Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /seller/login [post]
func (h *SellerAuthHandler) Login(c fiber.Ctx) error {
	var req domain.LoginReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.authUsecase.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Verify user has seller role
	hasRole := false
	for _, r := range res.Roles {
		if r == "seller" || r == "admin" {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account does not have seller permissions"})
	}

	return c.JSON(res)
}

// Refresh godoc
// @Summary Seller Refresh Token
// @Description Refresh seller access token
// @Tags Seller Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /seller/refresh [post]
func (h *SellerAuthHandler) Refresh(c fiber.Ctx) error {
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
// @Summary Seller Logout
// @Description Revoke seller refresh token
// @Tags Seller Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} map[string]string
// @Router /seller/logout [post]
func (h *SellerAuthHandler) Logout(c fiber.Ctx) error {
	var req domain.RefreshTokenReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_ = h.authUsecase.Logout(req.RefreshToken)
	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}
