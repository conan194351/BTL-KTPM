package impl

import (
	"context"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) i.OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (r *orderRepositoryImpl) Create(ctx context.Context, order interface{}) error {
	return nil
}

func (r *orderRepositoryImpl) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	var order interface{}
	err := r.db.WithContext(ctx).First(&order, id).Error
	return order, err
}

func (r *orderRepositoryImpl) Update(ctx context.Context, order models.Order) error {
	return r.db.WithContext(ctx).
		Save(order).
		Select("*").
		Where("id = ? AND deleted_at IS NULL", order.ID).Error
}

func (r *orderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}
