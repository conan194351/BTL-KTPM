package routes

import (
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	cnf := config.GetConfig()
	gin.SetMode(cnf.App.GetMode())
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())
	v1 := r.Group("/api/v1")
	AddHealthCheckRouter(v1)
	return r
}
