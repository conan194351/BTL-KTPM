package impl

import (
	"context"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/utils"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) i.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) Create(ctx context.Context, user interface{}) error {
	userModel, err := utils.ConvertToStruct[models.User](user)
	if err != nil {
		return err
	}
	if err := u.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepositoryImpl) GetByID(ctx context.Context, id uint) (interface{}, error) {
	var user models.User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
