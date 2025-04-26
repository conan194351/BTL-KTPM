package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null"`
	TotalPrice float64     `gorm:"not null"`
	Status     OrderStatus `gorm:"default:'pending'"`
	ProductID  uint        `gorm:"not null"`
}

type OrderStatus string

const (
	Pending          OrderStatus = "PENDING"
	Verified         OrderStatus = "VERIFIED"
	Paid             OrderStatus = "PAID"
	InventoryUpdated OrderStatus = "INVENTORY_UPDATED"
	VerifyFailed     OrderStatus = "VERIFY_FAILED"
	PaymentFailed    OrderStatus = "PAYMENT_FAILED"
	InventoryFailed  OrderStatus = "INVENTORY_FAILED"
	Completed        OrderStatus = "COMPLETED"
	EmailFailed      OrderStatus = "EMAIL_FAILED"
	Failed           OrderStatus = "FAILED"
)
