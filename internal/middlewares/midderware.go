package middlewares

import (
	"github.com/conan194351/BTL-KTPM/internal/repository/i"
	"github.com/conan194351/BTL-KTPM/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Middleware struct {
	jwt      jwt.Service
	userRepo i.UserRepository
}

func NewMiddleware(jwt jwt.Service, userRepo i.UserRepository) Middleware {
	return Middleware{
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token1 := c.Request.Header.Get("Authorization")
		if token1 == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		arr := strings.Split(token1, " ")
		if len(arr) <= 1 {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		token := arr[1]
		claims, err := m.jwt.VerifyJWTToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID := claims["data"].(float64)
		user, err := m.userRepo.GetByID(c, uint(userID))
		if err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
