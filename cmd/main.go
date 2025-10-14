package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/routes"
	gemini "github.com/GazDuckington/go-gin/pkgs/genai"
	"github.com/GazDuckington/go-gin/pkgs/minio"
	"github.com/GazDuckington/go-gin/pkgs/qdrant"
	"github.com/unidoc/unipdf/v4/common/license"
)

func main() {
	cfg := config.LoadConfig()

	// connect DB (safe to skip in dev if you want; error out otherwise)
	if err := database.Connect(cfg); err != nil {
		cfg.Logger.Warnf("database connection failed, continuing without DB: %v", err)
	} else {
		cfg.Logger.Info("database connected successfully")
	}

	if err := qdrant.Init(cfg); err != nil {
		cfg.Logger.Warnf("Faiure initiating Qdrant client: %v", err)
	} else {
		cfg.Logger.Info("qdrant client initialized")
	}

	if err := minio.Init(cfg); err != nil {
		cfg.Logger.Warnf("Faiure initiating Minio client: %v", err)
	} else {
		cfg.Logger.Info("minio client initialized")
	}

	if err := gemini.Init(context.Background(), cfg); err != nil {
		cfg.Logger.Warnf("Faiure initiating Gemini client: %v", err)
	} else {
		cfg.Logger.Info("gemini client initialized")
	}

	err := license.SetMeteredKey(cfg.UnidocKey)
	if err != nil {
		cfg.Logger.Warnf("Faiure initiating unidoc license: %v", err)
	} else {
		cfg.Logger.Info("unidoc license accepted")
	}

	// NOTE: we manage schema with migrate CLI; DO NOT call AutoMigrate here in prod.
	// If you want to auto-migrate for quick dev, you can call it explicitly.

	r := routes.SetupRouter(cfg)
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	cfg.Logger.Infof("starting server on %s", addr)

	// graceful shutdown basic pattern
	go func() {
		if err := r.Run(addr); err != nil {
			cfg.Logger.Fatalf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cfg.Logger.Info("shutting down")
}
