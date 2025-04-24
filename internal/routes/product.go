package routes

import "github.com/gin-gonic/gin"

func AddProductRouter(r *gin.RouterGroup) {
	authRouter := r.Group("/products")
	{
		authRouter.GET("/")
	}
}
