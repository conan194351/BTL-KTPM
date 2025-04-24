package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null"`
	TotalPrice float64     `gorm:"not null"`
	Status     string      `gorm:"default:'pending'"` // pending, paid, shipped, done
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}
