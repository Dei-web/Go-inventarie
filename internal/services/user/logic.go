package user

import (
	"context"

	"github.com/Dei-web/Go-inventarie/internal/middleware/auth"
	"github.com/Dei-web/Go-inventarie/internal/models"
	"github.com/Dei-web/Go-inventarie/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) GetAll(ctx context.Context) ([]types.ResponseData, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]types.ResponseData, len(users))
	for i, u := range users {
		response[i] = types.ResponseData{ID: u.ID, Name: u.Name, Email: u.Email}
	}
	return response, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*types.ResponseData, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &types.ResponseData{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

func (s *Service) Create(ctx context.Context, req *types.UsersCreate) (*types.ResponseData, error) {
	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &types.ResponseData{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *types.UsersUpdate) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashed, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.Update(ctx, id, updates)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token: token,
		User: types.ResponseData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
