package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/reverny/kratos-mono/gen/go/api/inventory/v1"
	"github.com/reverny/kratos-mono/services/inventory/internal/biz"
)

type InventoryService struct {
	v1.UnimplementedInventoryServer

	uc *biz.InventoryUsecase
}

func NewInventoryService(uc *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{uc: uc}
}

func (s *InventoryService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
	return s.uc.CreateProduct(ctx, req)
}

func (s *InventoryService) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
	return s.uc.GetProduct(ctx, req)
}

func (s *InventoryService) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	return s.uc.ListProducts(ctx, req)
}

func (s *InventoryService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Product, error) {
	return s.uc.UpdateProduct(ctx, req)
}

func (s *InventoryService) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*emptypb.Empty, error) {
	return s.uc.DeleteProduct(ctx, req)
}

func (s *InventoryService) UpdateStock(ctx context.Context, req *v1.UpdateStockRequest) (*v1.Product, error) {
	return s.uc.UpdateStock(ctx, req)
}
