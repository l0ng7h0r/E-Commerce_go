package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
	"github.com/l0ng7h0r/ecommerce/pkg/phajay"
)

type UserPaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewUserPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *UserPaymentHandler {
	return &UserPaymentHandler{paymentUsecase: paymentUsecase}
}

// CreatePayment godoc
// @Summary Create Phajay Payment Link
// @Description Make payment for an existing order by supplying order_id. Returns Phajay Payment URL.
// @Tags Customer Payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreatePaymentReq true "Payment Request"
// @Success 201 {object} domain.PaymentResponse
// @Failure 400 {object} map[string]string
// @Router /user/payments [post]
func (h *UserPaymentHandler) CreatePayment(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req domain.CreatePaymentReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	paymentRes, err := h.paymentUsecase.CreatePayment(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(paymentRes)
}

// PhajayWebhook godoc
// @Summary Phajay Payment Webhook
// @Description Webhook callback endpoint invoked by Phajay Payment Gateway
// @Tags Customer Payments
// @Accept json
// @Produce json
// @Param request body phajay.WebhookPayload true "Webhook Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /webhooks/phajay [post]
func (h *UserPaymentHandler) PhajayWebhook(c fiber.Ctx) error {
	var payload phajay.WebhookPayload
	if err := c.Bind().Body(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid webhook payload"})
	}

	if err := h.paymentUsecase.HandleWebhook(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Webhook processed successfully"})
}

// GetPaymentByOrder godoc
// @Summary Get Payment Status
// @Description Get payment status for order by order ID
// @Tags Customer Payments
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} domain.Payment
// @Failure 404 {object} map[string]string
// @Router /user/payments/order/{orderId} [get]
func (h *UserPaymentHandler) GetPaymentByOrder(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orderID := c.Params("orderId")
	roles, _ := c.Locals("user_roles").([]string)

	payment, err := h.paymentUsecase.GetPaymentByOrderID(orderID, userID, roles)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payment)
}
