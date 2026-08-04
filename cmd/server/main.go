package main

import (
	"net/http"

	"github.com/bookshelf/monolith/internal/handler"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)

	http.ListenAndServe(":8080", mux)
}
