package service

import (
	"context"
	"sync"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/repository"
	gemini "github.com/GazDuckington/go-gin/pkgs/genai"
	"github.com/GazDuckington/go-gin/pkgs/qdrant"
)

// CVWorkerService manages background CV evaluations
type CVWorkerService struct {
	cfg    *config.Config
	repo   repository.CVRepository
	jobs   chan string
	status sync.Map // map[cvID]string
}

// NewCVWorkerService creates and starts the worker
func NewCVWorkerService(cfg *config.Config, repo repository.CVRepository) *CVWorkerService {
	s := &CVWorkerService{
		cfg:  cfg,
		jobs: make(chan string, 100),
		repo: repo,
	}

	go s.workerLoop()
	return s
}

// EnqueueCV adds a CV to the evaluation queue and returns status info
func (s *CVWorkerService) EnqueueCV(cvID string) dto.WorkerStatusResponse {
	s.status.Store(cvID, "queued")
	s.jobs <- cvID

	return dto.WorkerStatusResponse{
		ID:     cvID,
		Status: "queued",
	}
}
func (s *CVWorkerService) setState(cvID, state string, eval *dto.CVEvaluationResponse) {
	s.status.Store(cvID, dto.WorkerStatusResponse{
		Status: state,
		Eval:   eval,
	})
}

// GetStatus retrieves the current status for a given CV
func (s *CVWorkerService) GetStatus(cvID string) dto.WorkerStatusResponse {
	if val, ok := s.status.Load(cvID); ok {
		state := val.(dto.WorkerStatusResponse)
		return dto.WorkerStatusResponse{
			ID:     cvID,
			Status: state.Status,
			Eval:   state.Eval,
		}
	}
	return dto.WorkerStatusResponse{
		ID:     cvID,
		Status: "not_found",
	}
}

// workerLoop processes jobs asynchronously
func (s *CVWorkerService) workerLoop() {
	ctx := context.Background()
	for cvID := range s.jobs {
		s.setState(cvID, "processing", nil)

		cv, err := s.repo.GetCv(ctx, cvID)
		if err != nil {
			s.cfg.Logger.Warnf("[worker] error getting CV %s: %v", cvID, err)
			s.setState(cvID, "error", nil)
			continue
		}

		if cv == nil {
			s.cfg.Logger.Warnf("[worker] CV not found for ID %s", cvID)
			s.setState(cvID, "not_found", nil)
			continue
		}

		qcv, err := qdrant.GetFromQdrant(ctx, cv, s.cfg)
		if err != nil {
			s.cfg.Logger.Warnf("[worker] failed to fetch Qdrant data for CV %s: %v", cvID, err)
			s.setState(cvID, "error", nil)
			continue
		}

		result, err := s.evaluateWithGemini(qcv)
		if err != nil {
			s.cfg.Logger.Warnf("[worker] evaluation failed for CV %s: %v", cvID, err)
			s.setState(cvID, "failed", nil)
			continue
		}

		s.setState(cvID, "done", result)
	}
}

// evaluateWithGemini is a placeholder for the LLM call
func (s *CVWorkerService) evaluateWithGemini(cv *dto.CVResponse) (*dto.CVEvaluationResponse, error) {
	eval, err := gemini.EvaluateCV(context.Background(), cv)
	if err != nil {
		s.cfg.Logger.Warnf("[worker] evaluating via gemini failed: %v", err)
		return nil, err
	}
	s.cfg.Logger.Printf("[worker] evaluating CV: %s | Title: %s", cv.ID, cv.Title)
	return eval, nil
}
