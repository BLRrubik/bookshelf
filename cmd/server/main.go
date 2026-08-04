package main

import (
	"net/http"

	"github.com/bookshelf/monolith/internal/config"
	"github.com/bookshelf/monolith/internal/handler"
)

func main() {

	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)

	http.ListenAndServe(":"+cfg.Port, mux)
}
