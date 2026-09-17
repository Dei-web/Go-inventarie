package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Dei-web/Go-inventarie/internal/types"
	"github.com/Dei-web/Go-inventarie/pkg/httperr"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		httperr.Internal(w, "failed to fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req types.UsersCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		httperr.Unprocessable(w, "validation failed", err.Error())
		return
	}

	user, err := h.service.Create(r.Context(), &req)
	if err != nil {
		httperr.Internal(w, "failed to create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		httperr.BadRequest(w, "invalid user id")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrUserNotFound) {
		httperr.NotFound(w, "user not found")
		return
	}
	if err != nil {
		httperr.Internal(w, "failed to fetch user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		httperr.BadRequest(w, "invalid user id")
		return
	}

	var req types.UsersUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		httperr.Unprocessable(w, "validation failed", err.Error())
		return
	}

	if err := h.service.Update(r.Context(), id, &req); errors.Is(err, ErrUserNotFound) {
		httperr.NotFound(w, "user not found")
		return
	} else if err != nil {
		httperr.Internal(w, "failed to update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		httperr.BadRequest(w, "invalid user id")
		return
	}

	if err := h.service.Delete(r.Context(), id); errors.Is(err, ErrUserNotFound) {
		httperr.NotFound(w, "user not found")
		return
	} else if err != nil {
		httperr.Internal(w, "failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}
