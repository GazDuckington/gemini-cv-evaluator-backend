package controller

import (
	"net/http"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	cfg *config.Config
}

func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{cfg: cfg}
}

func (h *HealthController) HealthCheck(c *gin.Context) {
	h.cfg.Logger.Info("Health check endpoint called")
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
