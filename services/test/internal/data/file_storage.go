package data

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/reverny/kratos-mono/services/test/internal/biz"
)

// LocalFileStorage implements FileStorage for local/development use
// In production, replace with MinIOStorage or S3Storage
type LocalFileStorage struct {
	baseURL    string // e.g., "http://localhost:8002/files"
	uploadPath string // local directory path
}

func NewLocalFileStorage(baseURL, uploadPath string) *LocalFileStorage {
	return &LocalFileStorage{
		baseURL:    baseURL,
		uploadPath: uploadPath,
	}
}

func (s *LocalFileStorage) GeneratePresignedURL(ctx context.Context, fileID, fileName, contentType string, expiresIn time.Duration) (*biz.PresignedURLInfo, error) {
	// For local storage, we use a custom upload endpoint
	// In production with S3/MinIO, this would generate real presigned URL
	
	uploadURL := fmt.Sprintf("%s/upload/%s", s.baseURL, fileID)
	downloadURL := fmt.Sprintf("%s/%s/%s", s.baseURL, fileID, fileName)
	
	return &biz.PresignedURLInfo{
		UploadURL:   uploadURL,
		DownloadURL: downloadURL,
		Method:      "PUT",
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		ExpiresIn: int64(expiresIn.Seconds()),
	}, nil
}

func (s *LocalFileStorage) GetFileURL(fileID string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, fileID)
}

func (s *LocalFileStorage) ConfirmUpload(ctx context.Context, fileID string) error {
	// For local storage, just check if file exists
	// In production, this might verify with S3/MinIO
	
	// TODO: Implement actual file existence check
	return nil
}

// GetFilePath returns the local file system path
func (s *LocalFileStorage) GetFilePath(fileID, fileName string) string {
	return filepath.Join(s.uploadPath, fileID, fileName)
}
