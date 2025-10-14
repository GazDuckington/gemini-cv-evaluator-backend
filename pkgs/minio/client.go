package minio

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

// Init sets up the MinIO client using your app config
func Init(cfg *config.Config) error {
	endpoint := fmt.Sprintf("%s:%d", cfg.MinioHost, cfg.MinioPort) // adjust if not localhost
	accessKey := cfg.MinioUser
	secretKey := cfg.MinioPass

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // true if you use HTTPS
	})
	if err != nil {
		return fmt.Errorf("failed to init minio client: %w", err)
	}

	Client = client
	return nil
}

// EnsureBucket checks if a bucket exists, and creates it if not
func EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		if err := Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

// UploadPDF uploads a PDF file from disk to MinIO
func UploadPDF(ctx context.Context, bucketName, objectName, filePath string) (minio.UploadInfo, error) {
	if Client == nil {
		return minio.UploadInfo{}, fmt.Errorf("minio client not initialized")
	}

	if filepath.Ext(filePath) != ".pdf" {
		return minio.UploadInfo{}, fmt.Errorf("only PDF files are allowed")
	}

	info, err := Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to upload PDF: %w", err)
	}
	return info, nil
}

// GetPresignedURL generates a presigned GET URL for an object
func GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	reqParams := make(url.Values)
	presignedURL, err := Client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}
