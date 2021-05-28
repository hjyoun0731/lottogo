package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lottogo/api"

	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("server starting...")

	router := httprouter.New()

	router.GET("/", api.Index)
	router.GET("/random", api.Random)

	router.PUT("/upload", api.UploadFile)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe fail")
		panic(err)
	}
}
