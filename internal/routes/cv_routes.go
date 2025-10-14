package routes

import (
	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/controller"
	"github.com/GazDuckington/go-gin/internal/middleware"
	"github.com/GazDuckington/go-gin/internal/repository"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterCvRoutes(r *gin.Engine, cfg *config.Config) {
	cvRepo := repository.NewCVRepository(database.DB, cfg)
	cvSvc := service.NewCVService(cvRepo, cfg)
	cvWrk := service.NewCVWorkerService(cfg, cvRepo)
	cvCtrl := controller.NewCvController(cvSvc, cfg, cvWrk)

	g := r.Group("/cv")
	g.Use(middleware.AuthRequired([]byte(cfg.JWTSecret), cfg.Logger))
	{
		g.POST("", cvCtrl.SubmitCv)
		g.GET("/:id", cvCtrl.GetCv)
		g.POST("/:id", cvCtrl.EvaluateCv)
		g.GET("status/:id", cvCtrl.GetEvalStatus)
		g.GET("result/:id", cvCtrl.EvaluationResult)
	}
}
