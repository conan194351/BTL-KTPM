package routes

import (
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(r *gin.RouterGroup) {
	authRouter := r.Group("/auth")
	{
		authRouter.POST("/login")
	}
}
