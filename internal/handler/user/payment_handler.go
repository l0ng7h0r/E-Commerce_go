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
// @Summary Create Phajay Payment Link (redirect)
// @Description Legacy redirect-based payment – returns external Phajay URL.
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

// GenerateQR godoc
// @Summary Generate Direct QR Code for Payment
// @Description Generates a Phajay QR Code string (EMVCo) for in-app scanning. No redirect needed.
// @Tags Customer Payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.GenerateQRReq true "QR Generation Request"
// @Success 201 {object} domain.GenerateQRResponse
// @Failure 400 {object} map[string]string
// @Router /user/payments/generate-qr [post]
func (h *UserPaymentHandler) GenerateQR(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req domain.GenerateQRReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	qrRes, err := h.paymentUsecase.GenerateQR(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(qrRes)
}

// CheckPaymentStatus godoc
// @Summary Check QR Payment Status (poll)
// @Description Polls Phajay for live transaction status and syncs the DB. Used by frontend polling.
// @Tags Customer Payments
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} domain.PaymentStatusResponse
// @Failure 400 {object} map[string]string
// @Router /user/payments/order/{orderId}/status [get]
func (h *UserPaymentHandler) CheckPaymentStatus(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orderID := c.Params("orderId")
	roles, _ := c.Locals("user_roles").([]string)

	status, err := h.paymentUsecase.CheckPaymentStatus(orderID, userID, roles)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(status)
}

// PhajayWebhook godoc
// @Summary Phajay Payment Webhook
// @Description Webhook callback endpoint invoked by Phajay Payment Gateway when a transaction completes or fails.
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
// @Description Get payment record for order by order ID
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
