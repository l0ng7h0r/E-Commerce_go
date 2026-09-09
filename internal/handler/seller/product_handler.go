package seller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type SellerProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewSellerProductHandler(productUsecase *usecase.ProductUsecase) *SellerProductHandler {
	return &SellerProductHandler{productUsecase: productUsecase}
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory godoc
// @Summary Seller Create Category
// @Description Create a new product category as seller
// @Tags Seller Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body createCategoryReq true "Category Request"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Router /seller/categories [post]
func (h *SellerProductHandler) CreateCategory(c fiber.Ctx) error {
	var req createCategoryReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	cat, err := h.productUsecase.CreateCategory(req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(cat)
}

// CreateProduct godoc
// @Summary Create Product
// @Description Add a new product for seller store
// @Tags Seller Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CreateProductReq true "Create Product Request"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /seller/products [post]
func (h *SellerProductHandler) CreateProduct(c fiber.Ctx) error {
	sellerID := c.Locals("user_id").(string)
	var req domain.CreateProductReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := h.productUsecase.CreateProduct(sellerID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

// GetMyProducts godoc
// @Summary List Seller Products
// @Description List products created by logged in seller
// @Tags Seller Products
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Product
// @Router /seller/products [get]
func (h *SellerProductHandler) GetMyProducts(c fiber.Ctx) error {
	sellerID := c.Locals("user_id").(string)
	products, err := h.productUsecase.GetProductsBySellerID(sellerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// UpdateProduct godoc
// @Summary Update Product
// @Description Update existing seller product details
// @Tags Seller Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body domain.UpdateProductReq true "Update Product Request"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /seller/products/{id} [put]
func (h *SellerProductHandler) UpdateProduct(c fiber.Ctx) error {
	sellerID := c.Locals("user_id").(string)
	productID := c.Params("id")
	var req domain.UpdateProductReq
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := h.productUsecase.UpdateProduct(productID, sellerID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

// DeleteProduct godoc
// @Summary Delete Product
// @Description Delete a product owned by seller
// @Tags Seller Products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /seller/products/{id} [delete]
func (h *SellerProductHandler) DeleteProduct(c fiber.Ctx) error {
	sellerID := c.Locals("user_id").(string)
	productID := c.Params("id")

	if err := h.productUsecase.DeleteProduct(productID, sellerID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
