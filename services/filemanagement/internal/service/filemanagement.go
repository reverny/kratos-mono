package service

import (
	"context"

	pb "github.com/reverny/kratos-mono/gen/go/api/filemanagement/v1"
	"github.com/reverny/kratos-mono/services/filemanagement/internal/biz"
)

type FilemanagementService struct {
	pb.UnimplementedFilemanagementServer

	fileUploadUC *biz.FileUploadUseCase
}

func NewFilemanagementService(fileUploadUC *biz.FileUploadUseCase) *FilemanagementService {
	return &FilemanagementService{
		fileUploadUC: fileUploadUC,
	}
}

func (s *FilemanagementService) RequestUploadUrl(ctx context.Context, req *pb.RequestUploadUrlRequest) (*pb.RequestUploadUrlReply, error) {
	metadata, urlInfo, err := s.fileUploadUC.RequestUpload(
		ctx,
		req.FileName,
		req.ContentType,
		req.FileSize,
		req.Description,
	)
	if err != nil {
		return nil, err
	}

	return &pb.RequestUploadUrlReply{
		FileId:      metadata.FileID,
		UploadUrl:   urlInfo.UploadURL,
		DownloadUrl: urlInfo.DownloadURL,
		Method:      urlInfo.Method,
		Headers:     urlInfo.Headers,
		ExpiresIn:   urlInfo.ExpiresIn,
	}, nil
}

func (s *FilemanagementService) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadReply, error) {
	fileURL, err := s.fileUploadUC.ConfirmUpload(ctx, req.FileId)
	if err != nil {
		return &pb.ConfirmUploadReply{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ConfirmUploadReply{
		Success: true,
		FileUrl: fileURL,
		Message: "File upload confirmed successfully",
	}, nil
}

func (s *FilemanagementService) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest) (*pb.GetFileInfoReply, error) {
	metadata, err := s.fileUploadUC.GetFileInfo(ctx, req.FileId)
	if err != nil {
		return nil, err
	}

	return &pb.GetFileInfoReply{
		FileId:      metadata.FileID,
		FileName:    metadata.FileName,
		FileUrl:     metadata.FileURL,
		ContentType: metadata.ContentType,
		FileSize:    metadata.FileSize,
		Description: metadata.Description,
		CreatedAt:   metadata.UploadedAt.Unix(),
	}, nil
}

func (s *FilemanagementService) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileReply, error) {
	err := s.fileUploadUC.DeleteFile(ctx, req.FileId)
	if err != nil {
		return &pb.DeleteFileReply{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteFileReply{
		Success: true,
		Message: "File deleted successfully",
	}, nil
}
