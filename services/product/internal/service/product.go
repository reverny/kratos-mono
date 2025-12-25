package service

import (
	"context"

	pb "github.com/reverny/kratos-mono/gen/go/api/product/v1"
	"github.com/reverny/kratos-mono/services/product/internal/biz"
)

type ProductService struct {
	pb.UnimplementedProductServer

	uc *biz.ProductUseCase
}

func NewProductService(uc *biz.ProductUseCase) *ProductService {
	return &ProductService{uc: uc}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	item, err := s.uc.Create(ctx, &biz.Product{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductReply{
		Data: &pb.ProductItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	item, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductReply{
		Data: &pb.ProductItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	items, total, err := s.uc.List(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	
	pbItems := make([]*pb.ProductItem, len(items))
	for i, item := range items {
		pbItems[i] = &pb.ProductItem{
			Id:   item.ID,
			Name: item.Name,
		}
	}
	
	return &pb.ListProductReply{
		Items: pbItems,
		Total: int32(total),
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	item, err := s.uc.Update(ctx, &biz.Product{
		ID:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductReply{
		Data: &pb.ProductItem{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductReply{Success: true}, nil
}
