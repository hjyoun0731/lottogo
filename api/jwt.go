package api

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserId int
	jwt.StandardClaims
}

var ExpirationTime time.Duration = time.Minute * 10
var JwtKey []byte = []byte("whsdhktxodkey")

func generateToken(c echo.Context, id int) (string, error) {
	expTime := time.Now().Add(ExpirationTime)
	claim := &Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", fmt.Errorf("token sign fail")
	}

	return tokenString, nil
}
