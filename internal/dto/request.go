package dto

type OderRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
}

type OrderWorkflowInput struct {
	OrderId uint
	UserId  uint
}
