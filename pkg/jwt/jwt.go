package jwt

import (
	"errors"
	"fmt"

	"github.com/conan194351/BTL-KTPM/internal/config"
	"github.com/golang-jwt/jwt"
)

var (
	HmacSecret = []byte(config.GetConfig().Server.SecretKey)
)

type Service interface {
	NewJWTToken(data uint, time int64) (*jwt.Token, *string, error)
	VerifyJWTToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct {
}

func NewJWTService() Service {
	return &jwtService{}
}

func (j *jwtService) NewJWTToken(data uint, time int64) (*jwt.Token, *string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time,
		"data": data,
	})
	tokenString, err := token.SignedString(HmacSecret)
	if err != nil {
		fmt.Printf("Error when signed string token " + err.Error())
		return nil, nil, fmt.Errorf("Unexpected error when signed string token " + err.Error())
	}
	return token, &tokenString, nil
}

func (j *jwtService) VerifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf(fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
			return nil, errors.New("unexpected signing method")
		}
		return HmacSecret, nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired {
			return nil, v
		}

		return nil, fmt.Errorf("Unexpected error when parse token: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return claims, nil
}
