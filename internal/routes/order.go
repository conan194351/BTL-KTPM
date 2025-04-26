package routes

import (
	"github.com/conan194351/BTL-KTPM/internal/handlers"
	"github.com/gin-gonic/gin"
)

func AddOrderRouter(r *gin.RouterGroup, orderHandler *handlers.OrderHandlerImpl) {
	authRouter := r.Group("/orders")
	{
		authRouter.POST("/", orderHandler.Create)
		authRouter.GET("/", orderHandler.Test)
	}
}
