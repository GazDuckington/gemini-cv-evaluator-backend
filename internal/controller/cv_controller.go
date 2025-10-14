package controller

import (
	"net/http"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/middleware"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CVController struct {
	svc service.CVService
	cfg *config.Config
	wrk *service.CVWorkerService
}

func NewCvController(s service.CVService, cfg *config.Config, wrk *service.CVWorkerService) *CVController {
	return &CVController{svc: s, cfg: cfg, wrk: wrk}
}

func (ctrl *CVController) SubmitCv(c *gin.Context) {
	var req dto.SubmitCvRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		ctrl.cfg.Logger.Warnf("Create bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("authClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user claims not found"})
		return
	}

	userClaims := claims.(*middleware.Claims)
	req.UserID = userClaims.UserID

	submitted, err := ctrl.svc.SubmitCV(c, req)
	if err != nil {
		ctrl.cfg.Logger.Errorf("Error submitting cv: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, submitted)
}

func (ctrl *CVController) GetCv(c *gin.Context) {
	cvID := c.Param("id")
	if cvID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing CV ID"})
		return
	}

	cv, err := ctrl.svc.GetCv(c.Request.Context(), cvID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cv == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cv})
}

func (ctrl *CVController) EvaluateCv(c *gin.Context) {
	cvID := c.Param("id")
	if cvID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing CV ID"})
		return
	}

	status := ctrl.wrk.EnqueueCV(cvID)

	c.JSON(http.StatusOK, gin.H{"data": status})
}

func (ctrl *CVController) GetEvalStatus(c *gin.Context) {
	cvID := c.Param("id")
	if cvID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing CV ID"})
		return
	}

	status := ctrl.wrk.GetStatus(cvID)

	c.JSON(http.StatusOK, gin.H{"data": status})
}

func (ctrl *CVController) EvaluationResult(c *gin.Context) {
	cvID := c.Param("id")
	if cvID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing CV ID"})
		return
	}

	status := ctrl.wrk.GetStatus(cvID)

	c.JSON(http.StatusOK, gin.H{"data": status})
}
