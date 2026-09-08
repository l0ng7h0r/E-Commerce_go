package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type UserOrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewUserOrderHandler(orderUsecase *usecase.OrderUsecase) *UserOrderHandler {
	return &UserOrderHandler{orderUsecase: orderUsecase}
}

// CreateOrder godoc
// @Summary Create Order (Checkout)
// @Description Create a new order from items currently in cart
// @Tags Customer Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateOrderReq false "Order Request"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Router /user/orders [post]
func (h *UserOrderHandler) CreateOrder(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req domain.CreateOrderReq
	_ = c.Bind().Body(&req)

	order, err := h.orderUsecase.CreateOrder(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetMyOrders godoc
// @Summary View Order History
// @Description Get list of past orders created by logged-in user
// @Tags Customer Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Router /user/orders [get]
func (h *UserOrderHandler) GetMyOrders(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orders, err := h.orderUsecase.GetUserOrderHistory(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// GetOrderByID godoc
// @Summary Get Order Details
// @Description Get specific order details by ID
// @Tags Customer Orders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 404 {object} map[string]string
// @Router /user/orders/{id} [get]
func (h *UserOrderHandler) GetOrderByID(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orderID := c.Params("id")
	roles, _ := c.Locals("user_roles").([]string)

	order, err := h.orderUsecase.GetOrderByID(orderID, userID, roles)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}
