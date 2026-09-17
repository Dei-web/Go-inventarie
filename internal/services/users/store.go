package users

import (
	"github.com/Dei-web/Go-inventarie/internal/middleware/auth"
	"github.com/Dei-web/Go-inventarie/internal/types"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUsers() ([]types.ResponseData, error) {
	var users []types.Users

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	response := make([]types.ResponseData, 0, len(users))

	for _, user := range users {
		response = append(response, types.ResponseData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return response, nil
}

func (s *Store) CreateUser(data *types.UsersCreate) (*types.UsersCreate, error) {
	passwordhashing, err := auth.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := &types.UsersCreate{
		Name:     data.Name,
		Email:    data.Email,
		Password: passwordhashing,
	}

	result := s.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *Store) DeleteUser(id int32) error {
	result := s.db.Delete(&types.Users{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store) UpdateUser(id int32, data *types.UsersUpdate) error {
	result := s.db.Model(&types.Users{}).Where("id=?", id).Updates(data)
	return result.Error
}

func (s *Store) GetUserByID(id int32) (error, *types.ResponseData) {
	var user types.Users

	result := s.db.First(&user, id)

	if result.Error != nil {
		return result.Error, nil
	}

	response := &types.ResponseData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return nil, response
}
