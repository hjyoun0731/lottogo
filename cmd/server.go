package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lottogo/api"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())

	e.GET("/", api.Index)
	e.GET("/random", api.Random)

	e.PUT("/upload", api.UploadFile)

	e.GET("/userinfo/get/:name", api.GetUserInfo)
	e.POST("/userinfo/signup", api.NewUserInfo)
	e.POST("/userinfo/signin", api.SignIn)

	e.GET("/download", api.DownloadFile)

	e.POST("/table/create", api.CreateTable)

	e.GET("/get/lotto/:round", api.GetLottoNum)

	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey:  []byte("secret"),
	// 	TokenLookup: "query:token",
	// }))

	e.Logger.Fatal(e.Start(":8080"))
}
