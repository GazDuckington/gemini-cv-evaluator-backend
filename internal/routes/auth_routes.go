package routes

import (
	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/controller"
	"github.com/GazDuckington/go-gin/internal/repository"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, cfg *config.Config) {
	repo := repository.NewUserRepository(database.DB, cfg)
	authService := service.NewAuthService(repo, cfg)
	authController := controller.NewAuthController(*authService, cfg)

	r.POST("/login", authController.Login)
}
