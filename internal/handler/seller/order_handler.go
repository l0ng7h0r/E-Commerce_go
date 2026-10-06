package seller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type SellerOrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewSellerOrderHandler(orderUsecase *usecase.OrderUsecase) *SellerOrderHandler {
	return &SellerOrderHandler{orderUsecase: orderUsecase}
}

// GetMyOrders godoc
// @Summary List Seller Store Orders
// @Description Get customer orders that contain products belonging to this seller
// @Tags Seller Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Router /seller/orders [get]
func (h *SellerOrderHandler) GetMyOrders(c fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(string)
	if !ok || sellerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	orders, err := h.orderUsecase.GetSellerOrders(sellerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if orders == nil {
		orders = []*domain.Order{}
	}
	return c.JSON(orders)
}

// UpdateOrderStatus godoc
// @Summary Update Order Status by Seller
// @Description Update order status (processing, shipped, completed, cancelled) for an order containing seller's products
// @Tags Seller Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body domain.UpdateOrderStatusReq true "Update Status Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /seller/orders/{id}/status [put]
func (h *SellerOrderHandler) UpdateOrderStatus(c fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(string)
	if !ok || sellerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	orderID := c.Params("id")
	var req domain.UpdateOrderStatusReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.orderUsecase.UpdateOrderStatusBySeller(orderID, sellerID, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order status updated successfully"})
}
