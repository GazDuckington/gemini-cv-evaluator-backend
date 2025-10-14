package routes

import (
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/controller"
	"github.com/GazDuckington/go-gin/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(
		middleware.RequestLogger(cfg.Logger),
		gin.Recovery(), // built-in panic recovery
	)

	// health
	healthCtrl := controller.NewHealthController(cfg)
	r.GET("/health", healthCtrl.HealthCheck)

	// NOTE: register domains
	RegisterUserRoutes(r, cfg)
	RegisterAuthRoutes(r, cfg)
	RegisterCvRoutes(r, cfg)
	return r
}
