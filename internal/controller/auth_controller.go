package controller

import (
	"net/http"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthController creates a new AuthController instance
func NewAuthController(authService service.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: &authService,
		cfg:         cfg,
	}
}

// Login handles POST /login
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.cfg.Logger.Warnf("invalid login request: %v", err)
		dto.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	ctx := c.Request.Context()

	token, err := ac.authService.Login(ctx, &req)
	if err != nil {
		ac.cfg.Logger.Warnf("login failed for username=%s: %v", req.Username, err)
		dto.Error(c, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	dto.Success(c, http.StatusOK, token, "login successful")
}
