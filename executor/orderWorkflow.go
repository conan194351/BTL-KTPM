package executor

import (
	"fmt"
	"github.com/conan194351/BTL-KTPM/internal/dto"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/pkg/logger"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type OrderWorkflow struct {
	logger   logger.Logger
	activity *Activities
}

func NewOrderWorkflow(
	a *Activities,
) *OrderWorkflow {
	return &OrderWorkflow{
		activity: a,
		logger:   logger.NewZapLogger("OrderWorkflow", true),
	}
}

func (o *OrderWorkflow) OrderWorkflow(ctx workflow.Context, input dto.OrderWorkflowInput) (*dto.OrderWorkflowResult, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	o.logger.Info(fmt.Sprintf("Starting order workflow for userID %s and orderID %s", input.UserId, input.OrderId), nil)

	// Step 1: Authenticate order
	var verifyResult bool
	err := workflow.ExecuteActivity(ctx, o.activity.VerifyOrderActivity, input.OrderId).Get(ctx, &verifyResult)
	if !verifyResult && err != nil {
		o.logger.Info("Order verification failed", nil)

		err = workflow.ExecuteActivity(ctx, o.activity.UpdateOrderStatusActivity, input.OrderId, models.VerifyFailed).Get(ctx, nil)
		if err != nil {
			o.logger.Error(err, "UpdateOrderStatusActivity failed", nil)
		}

		return &dto.OrderWorkflowResult{
			OrderID:     input.OrderId,
			OrderState:  string(models.VerifyFailed),
			ProcessedAt: workflow.Now(ctx),
		}, nil
	}

	// Step 2: Sending mail to confirm
	err = workflow.ExecuteActivity(ctx, o.activity.SendOrderConfirmationEmail, input.OrderId, input.UserId).Get(ctx, nil)
	if err != nil {
		o.logger.Error(err, "SendOrderConfirmationEmail failed", nil)
		err = workflow.ExecuteActivity(ctx, o.activity.UpdateOrderStatusActivity, input.OrderId, models.EmailFailed).Get(ctx, nil)
		if err != nil {
			o.logger.Error(err, "UpdateOrderStatusActivity failed", nil)
		}

		return &dto.OrderWorkflowResult{
			OrderID:     input.OrderId,
			OrderState:  string(models.EmailFailed),
			ProcessedAt: workflow.Now(ctx),
		}, nil
	}

	// Step 3: Process order payment
	err = workflow.ExecuteActivity(ctx, o.activity.ProcessPaymentActivity, input.OrderId).Get(ctx, nil)
	if err != nil {
		o.logger.Error(err, "ProcessPaymentActivity failed", nil)
		err = workflow.ExecuteActivity(ctx, o.activity.UpdateOrderStatusActivity, input.OrderId, models.PaymentFailed).Get(ctx, nil)
		if err != nil {
			o.logger.Error(err, "UpdateOrderStatusActivity failed", nil)
		}

		return &dto.OrderWorkflowResult{
			OrderID:     input.OrderId,
			OrderState:  string(models.PaymentFailed),
			ProcessedAt: workflow.Now(ctx),
		}, nil
	}

	// Step 4: Update inventory
	err = workflow.ExecuteActivity(ctx, o.activity.UpdateInventoryActivity, input.OrderId).Get(ctx, nil)
	if err != nil {
		o.logger.Error(err, "UpdateInventoryActivity failed", nil)
		err = workflow.ExecuteActivity(ctx, o.activity.UpdateOrderStatusActivity, input.OrderId, models.InventoryFailed).Get(ctx, nil)
		if err != nil {
			o.logger.Error(err, "UpdateOrderStatusActivity failed", nil)
		}

		return &dto.OrderWorkflowResult{
			OrderID:     input.OrderId,
			OrderState:  string(models.InventoryFailed),
			ProcessedAt: workflow.Now(ctx),
		}, nil
	}

	o.logger.Info("Order workflow completed successfully", nil)
	return &dto.OrderWorkflowResult{
		OrderID:     input.OrderId,
		OrderState:  string(models.Completed),
		ProcessedAt: workflow.Now(ctx),
	}, nil
}
