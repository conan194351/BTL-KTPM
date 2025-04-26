package i

import (
	"context"
	"github.com/conan194351/BTL-KTPM/internal/models"
)

type OrderRepository interface {
	Repository
	Update(ctx context.Context, order models.Order) error
	UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) error
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
}
