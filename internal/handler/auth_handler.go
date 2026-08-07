package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := h.services.UserService.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrInvalidUsername),
			errors.Is(err, service.ErrInvalidPassword),
			errors.Is(err, service.ErrUserExists),
			errors.Is(err, service.ErrUsernameExists):
			writeError(w, r, http.StatusBadRequest, err.Error())
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := h.services.UserService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, r, http.StatusBadRequest, "invalid credentials")

			return
		}

		writeError(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	user, err := h.services.UserService.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	writeJSON(w, http.StatusOK, user.ToPublic())
}

func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	user, err := h.services.UserService.Update(r.Context(), getUserID(r.Context()), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUsername),
			errors.Is(err, service.ErrUsernameExists):
			writeError(w, r, http.StatusBadRequest, "invalid username")
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	writeJSON(w, http.StatusOK, user.ToPublic())
}
