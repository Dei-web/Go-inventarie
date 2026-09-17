package user

import (
	"context"

	"github.com/Dei-web/Go-inventarie/internal/types"
)

type Repository interface {
	GetAll(ctx context.Context) ([]types.Users, error)
	GetByID(ctx context.Context, id int64) (*types.Users, error)
	Create(ctx context.Context, user *types.Users) error
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
}
