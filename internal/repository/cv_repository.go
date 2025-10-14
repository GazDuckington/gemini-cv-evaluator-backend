package repository

import (
	"context"

	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CVRepository interface {
	Submit(ctx context.Context, newCv *entity.CV) (*entity.CV, error)
	GetCv(ctx context.Context, id string) (*entity.CV, error)
}

type cvRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewCVRepository(db *gorm.DB, cfg *config.Config) CVRepository {
	return &cvRepository{
		db:     db,
		logger: cfg.Logger,
	}
}

func (r *cvRepository) Submit(ctx context.Context, newCv *entity.CV) (*entity.CV, error) {
	err := database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.Create(newCv).Error
	})
	if err != nil {
		return nil, err
	}
	return newCv, nil
}

func (r *cvRepository) GetCv(ctx context.Context, id string) (*entity.CV, error) {
	var cv entity.CV

	err := database.RunInTransaction(ctx, r.db, r.logger, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).First(&cv, "id = ?", id).Error
	})
	if err != nil {
		return nil, err
	}

	return &cv, nil
}
