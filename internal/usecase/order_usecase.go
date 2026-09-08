package usecase

import (
	"errors"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
)

type OrderUsecase struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewOrderUsecase(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *OrderUsecase) CreateOrder(userID string, req *domain.CreateOrderReq) (*domain.Order, error) {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var totalAmount float64
	var orderItems []domain.OrderItem

	for _, item := range cart.Items {
		product, err := u.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + item.ProductID)
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := &domain.Order{
		UserID:         userID,
		TotalAmount:    totalAmount,
		Status:         "pending",
		LogisticBranch: req.LogisticBranch,
	}

	createdOrder, err := u.orderRepo.CreateOrder(order, orderItems)
	if err != nil {
		return nil, err
	}

	// Clear cart after order creation
	_ = u.cartRepo.ClearCart(cart.ID)

	return createdOrder, nil
}

func (u *OrderUsecase) GetOrderByID(id, currentUserID string, roles []string) (*domain.Order, error) {
	order, err := u.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	// Admin can view any order, user can only view their own
	isAdmin := false
	for _, r := range roles {
		if r == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin && order.UserID != currentUserID {
		return nil, errors.New("unauthorized to view this order")
	}

	return order, nil
}

func (u *OrderUsecase) GetUserOrderHistory(userID string) ([]*domain.Order, error) {
	return u.orderRepo.GetOrdersByUserID(userID)
}

func (u *OrderUsecase) GetAllOrders() ([]*domain.Order, error) {
	return u.orderRepo.GetAllOrders()
}

func (u *OrderUsecase) UpdateOrderStatus(id, status string) error {
	return u.orderRepo.UpdateOrderStatus(id, status)
}
