package admin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/internal/usecase"
)

type AdminCategoryHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewAdminCategoryHandler(productUsecase *usecase.ProductUsecase) *AdminCategoryHandler {
	return &AdminCategoryHandler{productUsecase: productUsecase}
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory godoc
// @Summary Create Category
// @Description Create a new product category
// @Tags Admin Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body createCategoryReq true "Category Request"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Router /admin/categories [post]
func (h *AdminCategoryHandler) CreateCategory(c fiber.Ctx) error {
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
