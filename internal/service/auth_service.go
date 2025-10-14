package service

import (
	"context"
	"errors"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/repository"
	"github.com/GazDuckington/go-gin/pkgs/auth"
)

type AuthService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

// Login validates credentials and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, lp *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByLogin(ctx, lp.Username)
	if err != nil {
		return nil, errors.New("invalid email/username or password")
	}

	if err := user.ComparePassword(lp.Password); err != nil {
		return nil, errors.New("invalid email/username or password")
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Email, user.Role, s.cfg)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := auth.GenerateRefreshToken(user.Email, s.cfg)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	var profile *dto.ProfileResponse
	if user.Profile != nil {
		profile = &dto.ProfileResponse{
			ID:        user.Profile.ID,
			FullName:  user.Profile.FullName,
			Bio:       user.Profile.Bio,
			Phone:     user.Profile.Phone,
			AvatarURL: user.Profile.AvatarURL,
		}
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24h
		User: dto.UserResponse{
			ID:      user.ID,
			Email:   user.Email,
			Role:    user.Role,
			Profile: profile,
		},
	}, nil
}
