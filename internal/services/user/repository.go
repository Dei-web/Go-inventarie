package user

import (
	"context"

	"github.com/Dei-web/Go-inventarie/internal/models"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.Users, error)
	GetByID(ctx context.Context, id int64) (*models.Users, error)
	Create(ctx context.Context, user *models.Users) error
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
}
