package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lottogo/api"
)

type Stats struct {
	Uptime  time.Time     `json:"uptime"`
	ExpTime time.Duration `json:"exptime"`
}

func NewStats() *Stats {
	return &Stats{
		Uptime: time.Now(),
	}
}

func (s *Stats) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		s.ExpTime = time.Since(s.Uptime)
		log.Println(c.Path(), s.ExpTime)
		return err
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	s := NewStats()
	e.Use(s.Middleware)

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
