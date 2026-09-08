package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type UserCartHandler struct {
	cartUsecase *usecase.CartUsecase
}

func NewUserCartHandler(cartUsecase *usecase.CartUsecase) *UserCartHandler {
	return &UserCartHandler{cartUsecase: cartUsecase}
}

// GetCart godoc
// @Summary View Cart
// @Description View user's shopping cart
// @Tags Customer Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.Cart
// @Router /user/cart [get]
func (h *UserCartHandler) GetCart(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	cart, err := h.cartUsecase.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

// AddItem godoc
// @Summary Add Item to Cart
// @Description Add product to shopping cart
// @Tags Customer Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.AddCartItemReq true "Add Item Request"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]string
// @Router /user/cart/items [post]
func (h *UserCartHandler) AddItem(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req domain.AddCartItemReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	cart, err := h.cartUsecase.AddItem(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

// UpdateItem godoc
// @Summary Update Cart Item Quantity
// @Description Update quantity of item in cart
// @Tags Customer Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param request body domain.UpdateCartItemReq true "Update Item Request"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]string
// @Router /user/cart/items/{productId} [put]
func (h *UserCartHandler) UpdateItem(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("productId")
	var req domain.UpdateCartItemReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	cart, err := h.cartUsecase.UpdateItem(userID, productID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

// RemoveItem godoc
// @Summary Remove Item from Cart
// @Description Remove specific product from cart
// @Tags Customer Cart
// @Security BearerAuth
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} domain.Cart
// @Router /user/cart/items/{productId} [delete]
func (h *UserCartHandler) RemoveItem(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("productId")

	cart, err := h.cartUsecase.RemoveItem(userID, productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

// ClearCart godoc
// @Summary Clear Cart
// @Description Clear all items from cart
// @Tags Customer Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /user/cart [delete]
func (h *UserCartHandler) ClearCart(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	if err := h.cartUsecase.ClearCart(userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Cart cleared successfully"})
}
