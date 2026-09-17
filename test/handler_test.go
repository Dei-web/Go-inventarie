package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dei-web/Go-inventarie/internal/types"
	"github.com/Dei-web/Go-inventarie/internal/services/user"
)

func setupHandler() (*user.Handler, *MockRepository) {
	repo := NewMockRepository()
	svc := user.NewService(repo)
	handler := user.NewHandler(svc)
	return handler, repo
}

func TestHandlerCreateUser(t *testing.T) {
	h, _ := setupHandler()

	body := types.UsersCreate{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp types.ResponseData
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Name != body.Name {
		t.Errorf("expected name %q, got %q", body.Name, resp.Name)
	}

	if resp.ID == 0 {
		t.Error("expected user ID to be set")
	}
}

func TestHandlerCreateUserInvalidBody(t *testing.T) {
	h, _ := setupHandler()

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerCreateUserValidationFails(t *testing.T) {
	h, _ := setupHandler()

	body := types.UsersCreate{
		Name:     "J",
		Email:    "not-an-email",
		Password: "short",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestHandlerGetUsers(t *testing.T) {
	h, repo := setupHandler()

	repo.Create(context.Background(), &types.Users{Name: "User 1", Email: "u1@test.com", Password: "hashed"})

	svc := user.NewService(repo)
	svc.Create(context.Background(), &types.UsersCreate{Name: "User 2", Email: "u2@test.com", Password: "pass12345"})

	h = user.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlerDeleteUserNotFound(t *testing.T) {
	h, _ := setupHandler()

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandlerDeleteUserInvalidID(t *testing.T) {
	h, _ := setupHandler()

	req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerUpdateUser(t *testing.T) {
	h, repo := setupHandler()

	svc := user.NewService(repo)
	created, _ := svc.Create(context.Background(), &types.UsersCreate{Name: "John", Email: "john@test.com", Password: "pass12345"})
	h = user.NewHandler(svc)

	body := types.UsersUpdate{Name: "Jane"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonBody))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	updated, _ := svc.GetByID(context.Background(), int64(created.ID))
	if updated.Name != "Jane" {
		t.Errorf("expected name %q, got %q", "Jane", updated.Name)
	}
}
