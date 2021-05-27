package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Number struct {
	Num1 int
}

func main() {
	fmt.Println("server starting...")

	router := httprouter.New()

	router.GET("/", Index)
	router.GET("/random", Random)

	router.PUT("/upload", UploadFile)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe fail")
		panic(err)
	}
}
