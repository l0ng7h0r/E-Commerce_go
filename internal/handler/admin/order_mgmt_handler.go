package admin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type AdminOrderMgmtHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewAdminOrderMgmtHandler(orderUsecase *usecase.OrderUsecase) *AdminOrderMgmtHandler {
	return &AdminOrderMgmtHandler{orderUsecase: orderUsecase}
}

// GetAllOrders godoc
// @Summary List All System Orders
// @Description Get list of all orders across the system
// @Tags Admin Order Management
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Router /admin/orders [get]
func (h *AdminOrderMgmtHandler) GetAllOrders(c fiber.Ctx) error {
	orders, err := h.orderUsecase.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// UpdateOrderStatus godoc
// @Summary Update Order Status
// @Description Update processing status of an order
// @Tags Admin Order Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body domain.UpdateOrderStatusReq true "Update Status Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/orders/{id}/status [patch]
func (h *AdminOrderMgmtHandler) UpdateOrderStatus(c fiber.Ctx) error {
	orderID := c.Params("id")
	var req domain.UpdateOrderStatusReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.orderUsecase.UpdateOrderStatus(orderID, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order status updated successfully"})
}
