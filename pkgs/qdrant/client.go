package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/GazDuckington/go-gin/pkgs/minio"
	"github.com/qdrant/go-client/qdrant"
)

var QdrantClient *qdrant.Client

func Init(cfg *config.Config) error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.QdrantHost,
		Port:   cfg.QdrantGRPCPort,
		APIKey: cfg.QdrantKey,
		UseTLS: false,
	})
	if err != nil {
		return err
	}

	client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: cfg.MinioBucket,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     768,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	QdrantClient = client
	return nil
}

func NewCVPoint(cv *entity.CV) *qdrant.PointStruct {
	return &qdrant.PointStruct{
		Id:      qdrant.NewIDUUID(cv.ID),
		Vectors: qdrant.NewVectors(cv.Embedding...),
		Payload: qdrant.NewValueMap(map[string]any{
			"user_id": cv.UserID,
			"title":   cv.Title,
			"summary": cv.Summary,
			"file":    cv.FilePath,
		}),
	}
}

func StoreToQdrant(cv *entity.CV, cfg *config.Config) error {
	point := NewCVPoint(cv)

	_, err := QdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: cfg.MinioBucket,
		Points:         []*qdrant.PointStruct{point},
	})
	return err
}

func GetFromQdrant(ctx context.Context, cv *entity.CV, cfg *config.Config) (*dto.CVResponse, error) {
	// Fetch the point by ID
	resp, err := QdrantClient.Get(ctx, &qdrant.GetPoints{
		CollectionName: cfg.MinioBucket,
		Ids: []*qdrant.PointId{
			qdrant.NewIDUUID(cv.ID), // or NewIDUUID if your ID is a UUID
		},
		WithPayload: qdrant.NewWithPayload(true),
		WithVectors: qdrant.NewWithVectors(true),
	})
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("cv not found in qdrant for id: %s", cv.ID)
	}

	point := resp[0]
	payload := point.GetPayload()

	var file_path string

	fp, err := minio.GetPresignedURL(ctx, cfg.MinioBucket, payload["file"].GetStringValue(), 1*time.Hour)
	if err != nil {
		cfg.Logger.Warnf("[qdrant] failed to get presignedurl: %v", err)
		// fallback to the original file path if presign fails
		file_path = payload["file"].GetStringValue()
	} else {
		// use presigned URL if success
		file_path = fp
	}

	// Map payload back into entity.CV
	resCv := &dto.CVResponse{
		ID:       point.Id.GetUuid(), // helper below
		Title:    payload["title"].GetStringValue(),
		Summary:  payload["summary"].GetStringValue(),
		FilePath: file_path,
		UserID:   payload["user_id"].GetStringValue(),
	}

	// Convert vector back to []float32
	if v := point.Vectors; v != nil {
		resCv.Embedding = v.GetVector().Data
	}

	return resCv, nil
}
