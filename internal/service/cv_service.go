package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/GazDuckington/go-gin/internal/repository"
	gemini "github.com/GazDuckington/go-gin/pkgs/genai"
	"github.com/GazDuckington/go-gin/pkgs/minio"
	"github.com/GazDuckington/go-gin/pkgs/qdrant"
	"github.com/GazDuckington/go-gin/pkgs/utils"
)

type CVService interface {
	SubmitCV(ctx context.Context, req dto.SubmitCvRequest) (*entity.CV, error)
	GetCv(ctx context.Context, id string) (*dto.CVResponse, error)
}

type cvService struct {
	repo        repository.CVRepository
	minioBucket string
	cfg         *config.Config
}

func NewCVService(r repository.CVRepository, cfg *config.Config) CVService {
	return &cvService{
		repo:        r,
		minioBucket: cfg.MinioBucket,
		cfg:         cfg,
	}
}

func (s *cvService) SubmitCV(ctx context.Context, req dto.SubmitCvRequest) (*entity.CV, error) {
	// Ensure MinIO bucket exists
	if err := minio.EnsureBucket(ctx, s.minioBucket); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to ensure bucket: %w", err))
	}

	// Create unique file name
	objName := fmt.Sprintf("%d/%d_%s_%s", time.Now().Year(), time.Now().UnixNano(), req.UserID, req.File.Filename)

	// Upload file to MinIO
	uploadedFile, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open uploaded file: %w", err)
	}
	defer uploadedFile.Close()

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "cv-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("cannot create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up after upload

	// Copy uploaded content to temp file
	if _, err := io.Copy(tmpFile, uploadedFile); err != nil {
		return nil, fmt.Errorf("cannot copy uploaded file: %w", err)
	}

	// extract all texts from pdf
	text, err := utils.ExtractTextFromPDF(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %v", err)
	}
	// generate embedding from pdf
	embeds, err := gemini.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, err
	}

	uploaded, err := minio.UploadPDF(ctx, s.minioBucket, objName, tmpFile.Name())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to upload file: %w", err))
	}

	// Create entity
	newCv := &entity.CV{
		UserID:   req.UserID,
		Title:    req.Title,
		Summary:  text,
		FilePath: uploaded.Key,
	}

	// Save to DB
	embeddingJSON, err := json.Marshal(embeds)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding: %w", err)
	}
	mb, err := utils.BytesToFloat32sBinary(embeddingJSON)
	if err != nil {
		s.cfg.Logger.Warnf("error storing embeds: %v", err)
	}
	newCv.Embedding = mb
	created, err := s.repo.Submit(ctx, newCv)
	if err != nil {
		return nil, fmt.Errorf("failed to save CV: %w", err)
	}

	created.Embedding = embeds
	if err := qdrant.StoreToQdrant(created, s.cfg); err != nil {
		return nil, fmt.Errorf("qdrant upsert failed: %v", err)
	}

	return created, nil
}

func (s *cvService) GetCv(ctx context.Context, id string) (*dto.CVResponse, error) {
	cv, err := s.repo.GetCv(ctx, id)
	if err != nil {
		return nil, err
	}
	qcv, err := qdrant.GetFromQdrant(ctx, cv, s.cfg)
	if err != nil {
		return nil, err
	}
	// s.cfg.Logger.Debugf("summary and filepath:\n-%v\n-%v", cv.Summary, cv.FilePath)
	return qcv, nil
}
