package usecase

import (
	"errors"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
)

type CartUsecase struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartUsecase(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartUsecase {
	return &CartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *CartUsecase) GetCart(userID string) (*domain.Cart, error) {
	return u.cartRepo.GetOrCreateCart(userID)
}

func (u *CartUsecase) AddItem(userID string, req *domain.AddCartItemReq) (*domain.Cart, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	product, err := u.productRepo.GetProductByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if err := u.cartRepo.AddItem(cart.ID, req.ProductID, req.Quantity); err != nil {
		return nil, err
	}

	return u.cartRepo.GetOrCreateCart(userID)
}

func (u *CartUsecase) UpdateItem(userID, productID string, quantity int) (*domain.Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if err := u.cartRepo.UpdateItem(cart.ID, productID, quantity); err != nil {
		return nil, err
	}

	return u.cartRepo.GetOrCreateCart(userID)
}

func (u *CartUsecase) RemoveItem(userID, productID string) (*domain.Cart, error) {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if err := u.cartRepo.RemoveItem(cart.ID, productID); err != nil {
		return nil, err
	}

	return u.cartRepo.GetOrCreateCart(userID)
}

func (u *CartUsecase) ClearCart(userID string) error {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.ClearCart(cart.ID)
}
