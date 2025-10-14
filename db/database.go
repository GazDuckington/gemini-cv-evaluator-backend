package database

import (
	"time"

	"github.com/GazDuckington/go-gin/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := cfg.DatabaseDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// keep gorm logger minimal; we still use cfg.Logger for app logs
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	cfg.Logger.Infof("connected to DB %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	return nil
}
