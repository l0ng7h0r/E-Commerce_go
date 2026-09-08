package usecase

import (
	"errors"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) CreateCategory(name string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}
	return u.productRepo.CreateCategory(name)
}

func (u *ProductUsecase) GetAllCategories() ([]*domain.Category, error) {
	return u.productRepo.GetAllCategories()
}

func (u *ProductUsecase) GetCategoryByID(id string) (*domain.Category, error) {
	return u.productRepo.GetCategoryByID(id)
}

func (u *ProductUsecase) CreateProduct(sellerID string, req *domain.CreateProductReq) (*domain.Product, error) {
	product := &domain.Product{
		SellerID:    sellerID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}
	return u.productRepo.CreateProduct(product)
}

func (u *ProductUsecase) GetAllProducts() ([]*domain.Product, error) {
	return u.productRepo.GetAllProducts()
}

func (u *ProductUsecase) GetProductByID(id string) (*domain.Product, error) {
	return u.productRepo.GetProductByID(id)
}

func (u *ProductUsecase) GetProductsBySellerID(sellerID string) ([]*domain.Product, error) {
	return u.productRepo.GetProductsBySellerID(sellerID)
}

func (u *ProductUsecase) UpdateProduct(id, sellerID string, req *domain.UpdateProductReq) (*domain.Product, error) {
	p, err := u.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.Stock >= 0 {
		p.Stock = req.Stock
	}
	if req.CategoryID != "" {
		p.CategoryID = req.CategoryID
	}
	if req.ImageURL != "" {
		p.ImageURL = req.ImageURL
	}
	p.SellerID = sellerID

	if err := u.productRepo.UpdateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *ProductUsecase) DeleteProduct(id, sellerID string) error {
	return u.productRepo.DeleteProduct(id, sellerID)
}
