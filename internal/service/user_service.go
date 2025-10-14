package service

import (
	"context"

	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/GazDuckington/go-gin/internal/repository"
)

type UserService interface {
	GetAll(ctx context.Context) ([]dto.UserResponse, error)
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		var up *dto.ProfileResponse
		if u.Profile != nil { // assuming u.Profile is *model.Profile
			up = &dto.ProfileResponse{
				ID:        u.ID,
				Phone:     u.Profile.Phone,
				Bio:       u.Profile.Bio,
				AvatarURL: u.Profile.AvatarURL,
			}
		}
		out = append(out, dto.UserResponse{
			ID: u.ID, Email: u.Email, Role: u.Role, Profile: up,
		})
	}
	return out, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	resp := &dto.UserResponse{ID: u.ID, Email: u.Email}
	return resp, nil
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := &entity.User{Email: req.Email}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Email: user.Email}, nil
}
