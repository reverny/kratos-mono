package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/reverny/kratos-mono/gen/go/api/inventory/v1"
	"github.com/reverny/kratos-mono/services/inventory/internal/biz"
	"github.com/reverny/kratos-mono/services/inventory/internal/dto"
)

type InventoryService struct {
	v1.UnimplementedInventoryServer

	uc *biz.InventoryUsecase
}

func NewInventoryService(uc *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{uc: uc}
}

func (s *InventoryService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
	// Convert request to DTO
	createDTO := &dto.CreateProductDTO{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.Sku,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	// Call business logic with DTO
	productDTO, err := s.uc.CreateProduct(ctx, createDTO)
	if err != nil {
		return nil, err
	}

	// Convert DTO back to response
	return dtoToProto(productDTO), nil
}

func (s *InventoryService) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
	productDTO, err := s.uc.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return dtoToProto(productDTO), nil
}

func (s *InventoryService) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	query := &dto.ListProductsQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	products, total, err := s.uc.ListProducts(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert DTOs to proto
	protoProducts := make([]*v1.Product, len(products))
	for i, p := range products {
		protoProducts[i] = dtoToProto(p)
	}

	return &v1.ListProductsResponse{
		Products: protoProducts,
		Total:    total,
	}, nil
}

func (s *InventoryService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Product, error) {
	updateDTO := &dto.UpdateProductDTO{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	productDTO, err := s.uc.UpdateProduct(ctx, updateDTO)
	if err != nil {
		return nil, err
	}

	return dtoToProto(productDTO), nil
}

func (s *InventoryService) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *InventoryService) UpdateStock(ctx context.Context, req *v1.UpdateStockRequest) (*v1.Product, error) {
	stockDTO := &dto.UpdateStockDTO{
		ID:        req.Id,
		Quantity:  req.Quantity,
		Operation: req.Operation,
	}

	productDTO, err := s.uc.UpdateStock(ctx, stockDTO)
	if err != nil {
		return nil, err
	}

	return dtoToProto(productDTO), nil
}

// Helper function to convert DTO to proto
func dtoToProto(dto *dto.ProductDTO) *v1.Product {
	return &v1.Product{
		Id:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Sku:         dto.SKU,
		Price:       dto.Price,
		Stock:       dto.Stock,
		CreatedAt:   dto.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   dto.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
