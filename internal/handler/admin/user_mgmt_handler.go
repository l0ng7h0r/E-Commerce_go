package admin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type AdminUserMgmtHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAdminUserMgmtHandler(authUsecase *usecase.AuthUsecase) *AdminUserMgmtHandler {
	return &AdminUserMgmtHandler{authUsecase: authUsecase}
}

// CreateUser godoc
// @Summary Admin Create User
// @Description Create user with specific role (user, seller, admin)
// @Tags Admin User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.RegisterReq true "Create User Request"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /admin/users [post]
func (h *AdminUserMgmtHandler) CreateUser(c fiber.Ctx) error {
	var req domain.RegisterReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.authUsecase.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetAllUsers godoc
// @Summary List All Users
// @Description Get list of all registered users in system
// @Tags Admin User Management
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.User
// @Router /admin/users [get]
func (h *AdminUserMgmtHandler) GetAllUsers(c fiber.Ctx) error {
	users, err := h.authUsecase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GetUserByID godoc
// @Summary Get User Detail
// @Description Get detailed information of a specific user
// @Tags Admin User Management
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [get]
func (h *AdminUserMgmtHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.authUsecase.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete a user by ID
// @Tags Admin User Management
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/users/{id} [delete]
func (h *AdminUserMgmtHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.authUsecase.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
