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

// GetUsers godoc
// @Summary      Obtener todos los usuarios
// @Description  Retorna la lista de todos los usuarios
// @Tags         users
// @Produce      json
// @Success      200  {array}   types.ResponseData
// @Failure      500  {object}  httperr.HTTPError
// @Router       /users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		httperr.Internal(w, "failed to fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// CreateUser godoc
// @Summary      Crear un usuario
// @Description  Crea un nuevo usuario
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      types.UsersCreate  true  "Datos del usuario"
// @Success      201   {object}  types.ResponseData
// @Failure      400   {object}  httperr.HTTPError
// @Failure      422   {object}  httperr.HTTPError
// @Failure      500   {object}  httperr.HTTPError
// @Router       /users [post]
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

// GetUser godoc
// @Summary      Obtener un usuario
// @Description  Retorna un usuario por su ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID del usuario"
// @Success      200  {object}  types.ResponseData
// @Failure      400  {object}  httperr.HTTPError
// @Failure      404  {object}  httperr.HTTPError
// @Failure      500  {object}  httperr.HTTPError
// @Router       /users/{id} [get]
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

// UpdateUser godoc
// @Summary      Actualizar un usuario
// @Description  Actualiza un usuario por su ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "ID del usuario"
// @Param        user  body      types.UsersUpdate  true  "Datos a actualizar"
// @Success      204   "No Content"
// @Failure      400   {object}  httperr.HTTPError
// @Failure      404   {object}  httperr.HTTPError
// @Failure      422   {object}  httperr.HTTPError
// @Failure      500   {object}  httperr.HTTPError
// @Router       /users/{id} [put]
// @Security     ApiKeyAuth
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

// DeleteUser godoc
// @Summary      Eliminar un usuario
// @Description  Elimina un usuario por su ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID del usuario"
// @Success      204  "No Content"
// @Failure      400  {object}  httperr.HTTPError
// @Failure      404  {object}  httperr.HTTPError
// @Failure      500  {object}  httperr.HTTPError
// @Router       /users/{id} [delete]
// @Security     ApiKeyAuth
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

// Login godoc
// @Summary      Iniciar sesion
// @Description  Autentica un usuario y retorna un token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      types.LoginRequest  true  "Credenciales de acceso"
// @Success      200          {object}  types.LoginResponse
// @Failure      400          {object}  httperr.HTTPError
// @Failure      401          {object}  httperr.HTTPError
// @Failure      422          {object}  httperr.HTTPError
// @Router       /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		httperr.Unprocessable(w, "validation failed", err.Error())
		return
	}

	resp, err := h.service.Login(r.Context(), &req)
	if errors.Is(err, ErrInvalidCredentials) {
		httperr.Unauthorized(w, "invalid email or password")
		return
	}
	if err != nil {
		httperr.Internal(w, "failed to login")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func parseIDParam(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}
