package routes

import (
	"github.com/conan194351/BTL-KTPM/executor"
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/handlers"
	"github.com/conan194351/BTL-KTPM/internal/middlewares"
	"github.com/conan194351/BTL-KTPM/internal/repository/impl"
	"github.com/conan194351/BTL-KTPM/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"go.temporal.io/sdk/client"
	"log"
)

func InitRoutes(db *gorm.DB) *gin.Engine {
	cnf := config.GetConfig()
	gin.SetMode(cnf.App.GetMode())
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())

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

	//Service
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo, temporalClient, workflow)

	//Handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	v1 := r.Group("/api/v1")
	AddHealthCheckRouter(v1)

	AddAuthRouter(v1, db)
	AddProductRouter(v1, db)
	return r
}
