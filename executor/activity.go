package executor

import (
	"context"
	"fmt"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/mail"
	"github.com/conan194351/BTL-KTPM/pkg/utils"
	"time"
)

type Activities struct {
	orderRepo   i.OrderRepository
	productRepo i.ProductRepository
	userRepo    i.UserRepository
	mail        mail.MailService
}

func NewActivities(
	orderRepo i.OrderRepository,
	productRepo i.ProductRepository,
	userRepo i.UserRepository,
	mail mail.MailService,
) *Activities {
	return &Activities{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		mail:        mail,
	}
}

func (a *Activities) VerifyOrderActivity(ctx context.Context, orderID uint) (bool, error) {
	orderDetails, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return false, err
	}
	orderModel, err := utils.ConvertToStruct[models.Order](orderDetails)
	if err != nil {
		return false, err
	}
	if orderModel.Status != models.Pending {
		return false, fmt.Errorf("order has invalid status: %s", orderModel.Status)
	}

	item := orderModel.ProductID
	product, err := a.productRepo.GetByID(ctx, item)
	if err != nil {
		return false, err
	}
	productModel, err := utils.ConvertToStruct[models.Product](product)
	if err != nil {
		return false, err
	}
	if productModel.Stock < 1 {
		return false, fmt.Errorf("insufficient stock for product %d", item)
	}

	orderModel.Status = models.Verified
	orderModel.UpdatedAt = time.Now()
	if err := a.orderRepo.Update(ctx, orderModel); err != nil {
		return false, fmt.Errorf("failed to update order status: %w", err)
	}
	return true, nil
}

func (a *Activities) SendOrderConfirmationEmail(ctx context.Context, orderID uint, userID uint) error {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	userModel, err := utils.ConvertToStruct[models.User](user)
	if err != nil {
		return err
	}
	order, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	orderModel, err := utils.ConvertToStruct[models.Order](order)
	if err != nil {
		return fmt.Errorf("failed to convert order: %w", err)
	}
	product, err := a.productRepo.GetByID(ctx, orderModel.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	productModel, err := utils.ConvertToStruct[models.Product](product)
	if err != nil {
		return fmt.Errorf("failed to convert product: %w", err)
	}
	emailContent := fmt.Sprintf("Dear %s,\n\nYour order is a %s has been confirmed.\n\nThank you for your purchase!", userModel.Name, productModel.Name)
	err = a.mail.SendEmail(userModel.Email, "Order Confirmation", emailContent)
	if err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}

func (a *Activities) ProcessPaymentActivity(ctx context.Context, orderID uint) error {
	order, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	orderModel, err := utils.ConvertToStruct[models.Order](order)
	if err != nil {
		return fmt.Errorf("failed to convert order: %w", err)
	}
	// Mô phỏng xử lý thanh toán - trong thực tế sẽ gọi dịch vụ thanh toán
	// Thời gian chờ mô phỏng xử lý thanh toán
	time.Sleep(time.Second * 2)

	// Cập nhật trạng thái đơn hàng thành "payment_processed"
	orderModel.Status = models.Paid
	orderModel.UpdatedAt = time.Now()
	if err := a.orderRepo.Update(ctx, orderModel); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

func (a *Activities) UpdateInventoryActivity(ctx context.Context, orderID uint) error {
	order, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	orderModel, err := utils.ConvertToStruct[models.Order](order)
	if err != nil {
		return fmt.Errorf("failed to convert order: %w", err)
	}
	product, err := a.productRepo.GetByID(ctx, orderModel.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	productModel, err := utils.ConvertToStruct[models.Product](product)
	if err != nil {
		return fmt.Errorf("failed to convert product: %w", err)
	}

	if err := a.productRepo.UpdateStock(ctx, productModel.ID, productModel.Stock-1); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	orderModel.Status = models.InventoryUpdated
	orderModel.UpdatedAt = time.Now()
	if err := a.orderRepo.Update(ctx, orderModel); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (a *Activities) UpdateOrderStatusActivity(ctx context.Context, orderID uint, status models.OrderStatus) error {
	if err := a.orderRepo.UpdateOrderStatus(ctx, orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
