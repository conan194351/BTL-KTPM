package routes

import (
	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/conan194351/BTL-KTPM/internal/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB) *gin.Engine {
	cnf := config.GetConfig()
	gin.SetMode(cnf.App.GetMode())
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())
	v1 := r.Group("/api/v1")
	AddHealthCheckRouter(v1)

	AddAuthRouter(v1, db)
	return r
}
