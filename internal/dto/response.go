package dto

import (
	"github.com/conan194351/BTL-KTPM/internal/models"
	"time"
)

type OrderWorkflowResult struct {
	OrderID     uint
	OrderState  string
	ProcessedAt time.Time
}

type OrderResponse struct {
	ID         uint           `json:"id"`
	CreatedAt  string         `json:"created_at"`
	UserID     uint           `json:"user_id"`
	UserName   string         `json:"user_name"`
	UserEmail  string         `json:"user_email"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status"`
	OrderItems models.Product `json:"order_items"`
}
