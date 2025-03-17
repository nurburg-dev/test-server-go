package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the Go HTTP server!")
}

func main() {
	http.HandleFunc("/", helloHandler)

	fmt.Println("Starting server on port 9000...")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
