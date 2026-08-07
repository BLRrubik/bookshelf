package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
)

// Определяем собственный тип для ключей контекста
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

type Handler struct {
	services  *service.Service
	jwtSecret string
}

func New(services *service.Service, jwtSecret string) *Handler {
	return &Handler{
		services:  services,
		jwtSecret: jwtSecret,
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"ok", "version":"1.0.0", "timestamp":%d}`, time.Now().Unix())))
}

func Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf(`{"status":"ok", "version":"1.0.0", "timestamp":%d}, "checks":{"database": "ok"}`, time.Now().Unix()),
	))
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := domain.ErrorResponse{
		Code:      status,
		Message:   message,
		RequestID: r.Context().Value(requestIDKey).(string),
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}

func writeValidationError(w http.ResponseWriter, r *http.Request, details []domain.ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := domain.ErrorResponse{
		Code:      http.StatusBadRequest,
		Message:   "validation error",
		RequestID: r.Context().Value(requestIDKey).(string),
		Details:   details,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}
