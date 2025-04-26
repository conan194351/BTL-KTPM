package routes

import (
	"net/http"

	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProductRouter(r *gin.RouterGroup, db *gorm.DB) {
	productRouter := r.Group("/products")

	// GET /products - Lấy danh sách tất cả sản phẩm
	productRouter.GET("/", func(c *gin.Context) {
		products, err := models.GetAllProducts(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// POST /products - Tạo sản phẩm mới
	productRouter.POST("/", func(c *gin.Context) {
		var req struct {
			Name        string  `json:"name" binding:"required"`
			Description string  `json:"description"`
			Price       float64 `json:"price" binding:"required"`
			Stock       int     `json:"stock" binding:"required"`
			ImageURL    string  `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product := models.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			ImageURL:    req.ImageURL,
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	})
}
