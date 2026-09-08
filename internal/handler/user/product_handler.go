package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type UserProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewUserProductHandler(productUsecase *usecase.ProductUsecase) *UserProductHandler {
	return &UserProductHandler{productUsecase: productUsecase}
}

// GetAllProducts godoc
// @Summary List Products
// @Description Browse all available products
// @Tags Customer Products
// @Produce json
// @Success 200 {array} domain.Product
// @Router /user/products [get]
func (h *UserProductHandler) GetAllProducts(c fiber.Ctx) error {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// GetProductByID godoc
// @Summary Get Product Detail
// @Description Get product detail by ID
// @Tags Customer Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]string
// @Router /user/products/{id} [get]
func (h *UserProductHandler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

// GetAllCategories godoc
// @Summary List Categories
// @Description Get list of all product categories
// @Tags Customer Products
// @Produce json
// @Success 200 {array} domain.Category
// @Router /user/categories [get]
func (h *UserProductHandler) GetAllCategories(c fiber.Ctx) error {
	categories, err := h.productUsecase.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}
