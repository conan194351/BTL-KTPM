package routes

import (
	"github.com/conan194351/BTL-KTPM/executor"
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/handlers"
	"github.com/conan194351/BTL-KTPM/internal/middlewares"
	"github.com/conan194351/BTL-KTPM/internal/repository/impl"
	"github.com/conan194351/BTL-KTPM/internal/services"
	jwt2 "github.com/conan194351/BTL-KTPM/pkg/jwt"
	"github.com/conan194351/BTL-KTPM/pkg/mail"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
	"log"
)

func InitRoutes(db *gorm.DB) *gin.Engine {
	cnf := config.GetConfig()
	gin.SetMode(cnf.App.GetMode())
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())

	//DB
	config.InitDatabase()

	//Temporal
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233", // Địa chỉ của Temporal server
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}

	//Pkg
	mailService := mail.NewMailService()
	jwt := jwt2.NewJWTService()

	//Repo
	orderRepo := impl.NewOrderRepository(db)
	productRepo := impl.NewProductRepository(db)
	userRepo := impl.NewUserRepository(db)

	//Middleware
	midd := middlewares.NewMiddleware(jwt, userRepo)

	//Executor
	activities := executor.NewActivities(orderRepo, productRepo, userRepo, mailService)
	workflow := executor.NewOrderWorkflow(activities)

	//Service
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo, temporalClient, workflow, mailService)

	//Handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	v := r.Group("/api")
	AddAuthRouter(v, db)

	v1 := r.Group("/api/v1")
	v1.Use(midd.Auth())
	AddHealthCheckRouter(v1)
	AddProductRouter(v1, db)
	AddOrderRouter(v1, orderHandler)
	return r
}
