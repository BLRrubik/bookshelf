package main

import (
	"net/http"

	"github.com/bookshelf/monolith/internal/config"
	"github.com/bookshelf/monolith/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", handler.HealthHandler)

	http.ListenAndServe(":"+cfg.Port, r)
}
