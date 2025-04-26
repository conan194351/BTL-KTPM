package impl

import (
	"context"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/utils"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) i.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product interface{}) error {
	productModel, err := utils.ConvertToStruct[models.Product](product)
	if err != nil {
		return err
	}
	return r.db.Create(productModel).Error
}

func (r *productRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Update("stock", quantity).Error
}

func (r *productRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Update("status", status).Error
}
