package user

import (
	"context"
	"errors"

	"github.com/Dei-web/Go-inventarie/internal/models"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAll(ctx context.Context) ([]models.Users, error) {
	var users []models.Users
	if err := s.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*models.Users, error) {
	var user models.Users
	result := s.db.WithContext(ctx).First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*models.Users, error) {
	var user models.Users
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *Store) Create(ctx context.Context, user *models.Users) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *Store) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	result := s.db.WithContext(ctx).Model(&models.Users{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Delete(&models.Users{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
