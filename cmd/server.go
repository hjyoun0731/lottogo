package main

import (
	"github.com/labstack/echo"
	"github.com/lottogo/api"
)

func main() {
	e := echo.New()

	e.GET("/", api.Index)
	e.GET("/random", api.Random)

	e.PUT("/upload", api.UploadFile)

	e.GET("/userinfo/get", api.GetUserInfo)
	e.POST("/userinfo/signup", api.NewUserInfo)
	e.POST("/userinfo/signin", api.SignIn)

	// jc := middleware.JWTConfig{
	// 	SigningKey: []byte("secret"),
	// }
	// e.Use(middleware.JWTWithConfig(jc))

	e.Logger.Fatal(e.Start(":8080"))
}
