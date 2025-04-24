package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddHealthCheckRouter(r *gin.RouterGroup) {
	r.GET("healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
}
