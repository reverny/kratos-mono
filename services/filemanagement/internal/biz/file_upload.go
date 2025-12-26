package biz

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// FileUploadUseCase handles file upload business logic
type FileUploadUseCase struct {
	storage FileStorage
}

// FileStorage interface for storage backend (S3, MinIO, Local, etc.)
type FileStorage interface {
	// GeneratePresignedURL generates a presigned URL for uploading
	GeneratePresignedURL(ctx context.Context, fileID, fileName, contentType string, expiresIn time.Duration) (*PresignedURLInfo, error)
	
	// GetFileURL returns the public/download URL for a file
	GetFileURL(fileID string) string
	
	// ConfirmUpload verifies that a file was uploaded successfully (optional)
	ConfirmUpload(ctx context.Context, fileID string) error
	
	// DeleteFile deletes a file from storage
	DeleteFile(ctx context.Context, fileID string) error
	
	// GetFileInfo retrieves file metadata
	GetFileInfo(ctx context.Context, fileID string) (*FileMetadata, error)
}

// PresignedURLInfo contains presigned URL details
type PresignedURLInfo struct {
	UploadURL   string
	DownloadURL string
	Method      string
	Headers     map[string]string
	ExpiresIn   int64
}

// FileMetadata represents file metadata
type FileMetadata struct {
	FileID      string
	FileName    string
	ContentType string
	FileSize    int64
	Description string
	UploadedAt  time.Time
	FileURL     string
}

func NewFileUploadUseCase(storage FileStorage) *FileUploadUseCase {
	return &FileUploadUseCase{
		storage: storage,
	}
}

// RequestUpload generates presigned URL for file upload
func (uc *FileUploadUseCase) RequestUpload(ctx context.Context, fileName, contentType string, fileSize int64, description string) (*FileMetadata, *PresignedURLInfo, error) {
	// Generate unique file ID
	fileID := generateFileID()
	
	// Set expiration time (15 minutes)
	expiresIn := 15 * time.Minute
	
	// Generate presigned URL
	urlInfo, err := uc.storage.GeneratePresignedURL(ctx, fileID, fileName, contentType, expiresIn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	
	metadata := &FileMetadata{
		FileID:      fileID,
		FileName:    fileName,
		ContentType: contentType,
		FileSize:    fileSize,
		Description: description,
		UploadedAt:  time.Now(),
	}
	
	return metadata, urlInfo, nil
}

// ConfirmUpload confirms that file was uploaded successfully
func (uc *FileUploadUseCase) ConfirmUpload(ctx context.Context, fileID string) (string, error) {
	if err := uc.storage.ConfirmUpload(ctx, fileID); err != nil {
		return "", fmt.Errorf("failed to confirm upload: %w", err)
	}
	
	fileURL := uc.storage.GetFileURL(fileID)
	return fileURL, nil
}

// GetFileInfo returns file metadata
func (uc *FileUploadUseCase) GetFileInfo(ctx context.Context, fileID string) (*FileMetadata, error) {
	return uc.storage.GetFileInfo(ctx, fileID)
}

// DeleteFile deletes a file
func (uc *FileUploadUseCase) DeleteFile(ctx context.Context, fileID string) error {
	return uc.storage.DeleteFile(ctx, fileID)
}

// GetFileURL returns the download URL for a file
func (uc *FileUploadUseCase) GetFileURL(fileID string) string {
	return uc.storage.GetFileURL(fileID)
}

// generateFileID generates a unique file ID
func generateFileID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
