package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	ImageURL    string
}

func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Find(&products).Error
	return products, err
}
