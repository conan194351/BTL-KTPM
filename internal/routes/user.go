package routes

import (
	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	authRouter := r.Group("/users")
	{
		authRouter.GET("/me")
	}
}
