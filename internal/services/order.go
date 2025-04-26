package services

import (
	"context"
	"fmt"
	"github.com/conan194351/BTL-KTPM/executor"
	"github.com/conan194351/BTL-KTPM/internal/dto"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/mail"
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
	mailSrv       mail.MailService
}

func NewOrderService(
	orderRepo i.OrderRepository,
	userRepo i.UserRepository,
	productRepo i.ProductRepository,
	tempoClient client.Client,
	orderWorkflow *executor.OrderWorkflow,
	mailSrv mail.MailService,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo:     orderRepo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		tempoClient:   tempoClient,
		orderWorkflow: orderWorkflow,
		mailSrv:       mailSrv,
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

func (s *OrderServiceImpl) Test(ctx context.Context, orderID uint, userID uint) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	userModel, err := utils.ConvertToStruct[models.User](user)
	if err != nil {
		return err
	}
	emailContent := fmt.Sprintf("Dear %s,\n\nYour order with ID %d has been confirmed.\n\nThank you for your purchase!", userModel.Name, orderID)
	err = s.mailSrv.SendEmail(userModel.Email, "Order Confirmation", emailContent)
	if err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}
