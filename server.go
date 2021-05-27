package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("server starting...")
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", nil)
	fmt.Println("end")
}
