package test

import (
	"context"
	"testing"

	"github.com/Dei-web/Go-inventarie/internal/services/user"
	"github.com/Dei-web/Go-inventarie/internal/types"
)

const testJWTSecret = "test-secret-key"

func TestServiceCreateUser(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	req := &types.UsersCreate{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	resp, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Name != req.Name {
		t.Errorf("expected name %q, got %q", req.Name, resp.Name)
	}

	if resp.Email != req.Email {
		t.Errorf("expected email %q, got %q", req.Email, resp.Email)
	}

	if resp.ID == 0 {
		t.Error("expected user ID to be set")
	}
}

func TestServiceGetAll(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	svc.Create(context.Background(), &types.UsersCreate{
		Name: "User 1", Email: "user1@test.com", Password: "pass12345",
	})
	svc.Create(context.Background(), &types.UsersCreate{
		Name: "User 2", Email: "user2@test.com", Password: "pass12345",
	})

	users, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestServiceGetByID(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	created, _ := svc.Create(context.Background(), &types.UsersCreate{
		Name: "John", Email: "john@test.com", Password: "pass12345",
	})

	resp, err := svc.GetByID(context.Background(), int64(created.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Name != "John" {
		t.Errorf("expected name %q, got %q", "John", resp.Name)
	}
}

func TestServiceGetByIDNotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	_, err := svc.GetByID(context.Background(), 999)
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestServiceUpdate(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	created, _ := svc.Create(context.Background(), &types.UsersCreate{
		Name: "John", Email: "john@test.com", Password: "pass12345",
	})

	err := svc.Update(context.Background(), int64(created.ID), &types.UsersUpdate{
		Name: "Jane",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetByID(context.Background(), int64(created.ID))
	if updated.Name != "Jane" {
		t.Errorf("expected name %q, got %q", "Jane", updated.Name)
	}
}

func TestServiceDelete(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	created, _ := svc.Create(context.Background(), &types.UsersCreate{
		Name: "John", Email: "john@test.com", Password: "pass12345",
	})

	err := svc.Delete(context.Background(), int64(created.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.GetByID(context.Background(), int64(created.ID))
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	err := svc.Delete(context.Background(), 999)
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestServiceLogin(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	svc.Create(context.Background(), &types.UsersCreate{
		Name: "John", Email: "john@test.com", Password: "pass12345",
	})

	resp, err := svc.Login(context.Background(), &types.LoginRequest{
		Email:    "john@test.com",
		Password: "pass12345",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("expected token to be set")
	}

	if resp.User.Name != "John" {
		t.Errorf("expected name %q, got %q", "John", resp.User.Name)
	}
}

func TestServiceLoginInvalidEmail(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	_, err := svc.Login(context.Background(), &types.LoginRequest{
		Email:    "nonexistent@test.com",
		Password: "pass12345",
	})
	if err != user.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestServiceLoginInvalidPassword(t *testing.T) {
	repo := NewMockRepository()
	svc := user.NewService(repo, testJWTSecret)

	svc.Create(context.Background(), &types.UsersCreate{
		Name: "John", Email: "john@test.com", Password: "pass12345",
	})

	_, err := svc.Login(context.Background(), &types.LoginRequest{
		Email:    "john@test.com",
		Password: "wrongpassword",
	})
	if err != user.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}
