package test

import (
	"context"

	"github.com/Dei-web/Go-inventarie/internal/types"
	"github.com/Dei-web/Go-inventarie/internal/services/user"
)

type MockRepository struct {
	users     []types.Users
	nextID    uint
	getErr    error
	createErr error
	updateErr error
	deleteErr error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:  []types.Users{},
		nextID: 1,
	}
}

func (m *MockRepository) SetGetError(err error) {
	m.getErr = err
}

func (m *MockRepository) SetCreateError(err error) {
	m.createErr = err
}

func (m *MockRepository) SetUpdateError(err error) {
	m.updateErr = err
}

func (m *MockRepository) SetDeleteError(err error) {
	m.deleteErr = err
}

func (m *MockRepository) GetAll(ctx context.Context) ([]types.Users, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.users, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*types.Users, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, u := range m.users {
		if u.ID == uint(id) {
			return &u, nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (m *MockRepository) Create(ctx context.Context, user *types.Users) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = m.nextID
	m.nextID++
	m.users = append(m.users, *user)
	return nil
}

func (m *MockRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, u := range m.users {
		if u.ID == uint(id) {
			if name, ok := updates["name"].(string); ok && name != "" {
				m.users[i].Name = name
			}
			if email, ok := updates["email"].(string); ok && email != "" {
				m.users[i].Email = email
			}
			if password, ok := updates["password"].(string); ok && password != "" {
				m.users[i].Password = password
			}
			return nil
		}
	}
	return user.ErrUserNotFound
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, u := range m.users {
		if u.ID == uint(id) {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return user.ErrUserNotFound
}
