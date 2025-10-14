package repository

import (
	"context"
	"errors"

	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type userRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(db *gorm.DB, cfg *config.Config) UserRepository {
	return &userRepository{
		db:     db,
		logger: cfg.Logger,
	}
}

func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.Find(&users).Error
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.First(&user, id).Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User
	err := database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.Where("email = ?", login).
			Preload("Profile").
			First(&user).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}
