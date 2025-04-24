package routes

import "github.com/gin-gonic/gin"

func AddOrderRouter(r *gin.RouterGroup) {
	authRouter := r.Group("/orders")
	{
		authRouter.POST("/")
		authRouter.PUT("/:id")
	}
}
