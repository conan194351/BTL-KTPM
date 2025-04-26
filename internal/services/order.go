package services

import (
	"context"
	"fmt"
	"github.com/conan194351/BTL-KTPM/executor"
	"github.com/conan194351/BTL-KTPM/internal/dto"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/utils"
	"go.temporal.io/sdk/client"
	"time"
)

type OrderServiceImpl struct {
	orderRepo     i.OrderRepository
	userRepo      i.UserRepository
	productRepo   i.ProductRepository
	tempoClient   client.Client
	orderWorkflow *executor.OrderWorkflow
}

func NewOrderService(
	orderRepo i.OrderRepository,
	userRepo i.UserRepository,
	productRepo i.ProductRepository,
	tempoClient client.Client,
	orderWorkflow *executor.OrderWorkflow,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo:     orderRepo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		tempoClient:   tempoClient,
		orderWorkflow: orderWorkflow,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *dto.OderRequest) (*dto.OrderResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	userModel, err := utils.ConvertToStruct[models.User](user)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product %d: %w", req.ProductID, err)
	}
	productModel, err := utils.ConvertToStruct[models.Product](product)
	if err != nil {
		return nil, err
	}
	order := models.Order{
		UserID:     req.UserID,
		TotalPrice: productModel.Price,
		Status:     models.Pending,
		ProductID:  productModel.ID,
	}
	orderModel, err := s.orderRepo.CreateOrder(ctx, &order)
	if err != nil {
		return nil, err
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("order-workflow-%d", order.ID),
		TaskQueue: "order-processing-queue",
	}
	workflowInput := dto.OrderWorkflowInput{
		OrderId: order.ID,
		UserId:  req.UserID,
	}

	_, err = s.tempoClient.ExecuteWorkflow(ctx, workflowOptions, s.orderWorkflow.OrderWorkflow, workflowInput)
	if err != nil {
		return nil, fmt.Errorf("failed to start order workflow: %w", err)
	}
	response := &dto.OrderResponse{
		ID:         orderModel.ID,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UserID:     orderModel.UserID,
		UserName:   userModel.Name,
		UserEmail:  userModel.Email,
		TotalPrice: order.TotalPrice,
		Status:     string(orderModel.Status),
		OrderItems: productModel,
	}
	return response, nil
}

func (s *OrderServiceImpl) Test(ctx context.Context, orderID uint) (bool, error) {
	orderDetails, err := s.orderRepo.GetByID(ctx, orderID)
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
	product, err := s.productRepo.GetByID(ctx, item)
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
	if err := s.orderRepo.Update(ctx, orderModel); err != nil {
		return false, fmt.Errorf("failed to update order status: %w", err)
	}
	return true, nil
}
