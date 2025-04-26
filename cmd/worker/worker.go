package main

import (
	"github.com/conan194351/BTL-KTPM/executor"
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/repository/impl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	//DB
	config.InitDatabase()
	db := config.GetDB()

	//Temporal
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233", // Địa chỉ của Temporal server
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	//Repo
	orderRepo := impl.NewOrderRepository(db)
	productRepo := impl.NewProductRepository(db)
	userRepo := impl.NewUserRepository(db)

	//Executor
	activities := executor.NewActivities(orderRepo, productRepo, userRepo)
	workflow := executor.NewOrderWorkflow(activities)
	w := worker.New(temporalClient, "order-processing-queue", worker.Options{})
	w.RegisterWorkflow(workflow.OrderWorkflow)
	w.RegisterActivity(activities.VerifyOrderActivity)
	w.RegisterActivity(activities.ProcessPaymentActivity)
	w.RegisterActivity(activities.UpdateInventoryActivity)
	w.RegisterActivity(activities.UpdateOrderStatusActivity)
	log.Println("Starting Temporal worker...")
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Error running worker: %v", err)
	}
}
