package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type UserAuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewUserAuthHandler(authUsecase *usecase.AuthUsecase) *UserAuthHandler {
	return &UserAuthHandler{authUsecase: authUsecase}
}

// Register godoc
// @Summary Customer Registration
// @Description Register a new customer account
// @Tags Customer Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterReq true "Register Request"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /user/register [post]
func (h *UserAuthHandler) Register(c fiber.Ctx) error {
	var req domain.RegisterReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.Role = "user" // Ensure customer role

	res, err := h.authUsecase.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// Login godoc
// @Summary Customer Login
// @Description Login as customer
// @Tags Customer Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginReq true "Login Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /user/login [post]
func (h *UserAuthHandler) Login(c fiber.Ctx) error {
	var req domain.LoginReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.authUsecase.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// Refresh godoc
// @Summary Refresh Token
// @Description Refresh customer access token
// @Tags Customer Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /user/refresh [post]
func (h *UserAuthHandler) Refresh(c fiber.Ctx) error {
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
// @Summary Customer Logout
// @Description Revoke customer refresh token
// @Tags Customer Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Request"
// @Success 200 {object} map[string]string
// @Router /user/logout [post]
func (h *UserAuthHandler) Logout(c fiber.Ctx) error {
	var req domain.RefreshTokenReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_ = h.authUsecase.Logout(req.RefreshToken)
	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}
