package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/reverny/kratos-mono/services/filemanagement/internal/biz"
)

// LocalFileStorage implements FileStorage for local/development use
type LocalFileStorage struct {
	baseURL    string // e.g., "http://localhost:8005/files"
	uploadPath string // local directory path
}

func NewLocalFileStorage(baseURL, uploadPath string) *LocalFileStorage {
	// Create upload directory if it doesn't exist
	os.MkdirAll(uploadPath, 0755)
	
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
	// For local storage, check if file exists in the upload directory
	filePath := filepath.Join(s.uploadPath, fileID)
	
	// Check if directory exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", fileID)
	}
	
	return nil
}

func (s *LocalFileStorage) DeleteFile(ctx context.Context, fileID string) error {
	filePath := filepath.Join(s.uploadPath, fileID)
	
	// Remove the file directory
	if err := os.RemoveAll(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (s *LocalFileStorage) GetFileInfo(ctx context.Context, fileID string) (*biz.FileMetadata, error) {
	filePath := filepath.Join(s.uploadPath, fileID)
	
	// Check if directory exists
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	return &biz.FileMetadata{
		FileID:     fileID,
		FileURL:    s.GetFileURL(fileID),
		UploadedAt: info.ModTime(),
	}, nil
}

// GetFilePath returns the local file system path
func (s *LocalFileStorage) GetFilePath(fileID, fileName string) string {
	return filepath.Join(s.uploadPath, fileID, fileName)
}

// TODO: Implement MinIOStorage and S3Storage for production
// type MinIOStorage struct { ... }
// type S3Storage struct { ... }
