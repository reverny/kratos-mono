package service

import (
	"context"
	"fmt"

	pb "github.com/reverny/kratos-mono/gen/go/api/test/v1"
	"github.com/reverny/kratos-mono/services/test/internal/biz"
	"github.com/reverny/kratos-mono/services/test/internal/dto"
)

type TestService struct {
	pb.UnimplementedTestServer

	uc         *biz.TestUseCase
	fileUpload *biz.FileUploadUseCase
}

func NewTestService(uc *biz.TestUseCase, fileUpload *biz.FileUploadUseCase) *TestService {
	return &TestService{
		uc:         uc,
		fileUpload: fileUpload,
	}
}

func (s *TestService) CreateTest(ctx context.Context, req *pb.CreateTestRequest) (*pb.CreateTestReply, error) {
	item, err := s.uc.Create(ctx, &dto.CreateTestDTO{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTestReply{
		Data: &pb.TestItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *TestService) GetTest(ctx context.Context, req *pb.GetTestRequest) (*pb.GetTestReply, error) {
	item, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTestReply{
		Data: &pb.TestItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *TestService) ListTest(ctx context.Context, req *pb.ListTestRequest) (*pb.ListTestReply, error) {
	items, total, err := s.uc.List(ctx, &dto.ListTestQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	
	pbItems := make([]*pb.TestItem, len(items))
	for i, item := range items {
		pbItems[i] = &pb.TestItem{
			Id:   item.ID,
			Name: item.Name,
		}
	}
	
	return &pb.ListTestReply{
		Items: pbItems,
		Total: int32(total),
	}, nil
}

func (s *TestService) UpdateTest(ctx context.Context, req *pb.UpdateTestRequest) (*pb.UpdateTestReply, error) {
	item, err := s.uc.Update(ctx, &dto.UpdateTestDTO{
		ID:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateTestReply{
		Data: &pb.TestItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *TestService) DeleteTest(ctx context.Context, req *pb.DeleteTestRequest) (*pb.DeleteTestReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTestReply{Success: true}, nil
}

// RequestFileUpload generates presigned URL for file upload
func (s *TestService) RequestFileUpload(ctx context.Context, req *pb.RequestFileUploadRequest) (*pb.UploadUrlReply, error) {
	// Validate request
	if req.FileName == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if req.ContentType == "" {
		return nil, fmt.Errorf("content type is required")
	}
	
	// Request upload URL from use case
	metadata, urlInfo, err := s.fileUpload.RequestUpload(
		ctx,
		req.FileName,
		req.ContentType,
		req.FileSize,
		req.Description,
	)
	if err != nil {
		return nil, err
	}
	
	return &pb.UploadUrlReply{
		FileId:      metadata.FileID,
		UploadUrl:   urlInfo.UploadURL,
		DownloadUrl: urlInfo.DownloadURL,
		ExpiresIn:   urlInfo.ExpiresIn,
		Method:      urlInfo.Method,
		Headers:     urlInfo.Headers,
	}, nil
}

// ConfirmFileUpload confirms that file was uploaded successfully
func (s *TestService) ConfirmFileUpload(ctx context.Context, req *pb.ConfirmFileUploadRequest) (*pb.ConfirmFileUploadReply, error) {
	if req.FileId == "" {
		return nil, fmt.Errorf("file ID is required")
	}
	
	fileURL, err := s.fileUpload.ConfirmUpload(ctx, req.FileId)
	if err != nil {
		return nil, err
	}
	
	return &pb.ConfirmFileUploadReply{
		Success: true,
		FileUrl: fileURL,
		Message: "File upload confirmed successfully",
	}, nil
}
