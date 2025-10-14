package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunInTransaction(ctx context.Context, db *gorm.DB, logger *logrus.Logger, fn func(tx *gorm.DB) error) error {
	start := time.Now()
	logger.Debug("starting transaction")

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"duration": time.Since(start).String(),
			"error":    err.Error(),
		}).Error("transaction rolled back")
		return err
	}

	logger.WithFields(logrus.Fields{
		"duration": time.Since(start).String(),
	}).Debug("transaction committed")
	return nil
}
