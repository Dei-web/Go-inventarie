package httperr

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func New(w http.ResponseWriter, code int, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	New(w, http.StatusBadRequest, message, "")
}

func Unauthorized(w http.ResponseWriter, message string) {
	New(w, http.StatusUnauthorized, message, "")
}

func NotFound(w http.ResponseWriter, message string) {
	New(w, http.StatusNotFound, message, "")
}

func Internal(w http.ResponseWriter, message string) {
	New(w, http.StatusInternalServerError, message, "")
}

func Conflict(w http.ResponseWriter, message string) {
	New(w, http.StatusConflict, message, "")
}

func Unprocessable(w http.ResponseWriter, message string, details string) {
	New(w, http.StatusUnprocessableEntity, message, details)
}
