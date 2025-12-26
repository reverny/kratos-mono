package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/reverny/kratos-mono/services/inventory/internal/dto"
)

// InventoryRepo is a Inventory repo.
type InventoryRepo interface {
	CreateProduct(context.Context, *dto.CreateProductDTO) (*dto.ProductDTO, error)
	GetProduct(context.Context, string) (*dto.ProductDTO, error)
	ListProducts(context.Context, *dto.ListProductsQuery) ([]*dto.ProductDTO, int32, error)
	UpdateProduct(context.Context, *dto.UpdateProductDTO) (*dto.ProductDTO, error)
	DeleteProduct(context.Context, string) error
	UpdateStock(context.Context, *dto.UpdateStockDTO) (*dto.ProductDTO, error)
}

// InventoryUsecase is a Inventory usecase.
type InventoryUsecase struct {
	repo InventoryRepo
	log  *log.Helper
}

// NewInventoryUsecase new a Inventory usecase.
func NewInventoryUsecase(repo InventoryRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateProduct creates a Product.
func (uc *InventoryUsecase) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*dto.ProductDTO, error) {
	uc.log.WithContext(ctx).Infof("CreateProduct: %v", req.Name)
	
	// Business logic here (validation, business rules, etc.)
	if req.Stock < 0 {
		req.Stock = 0
	}
	
	return uc.repo.CreateProduct(ctx, req)
}

// GetProduct gets a Product by ID.
func (uc *InventoryUsecase) GetProduct(ctx context.Context, id string) (*dto.ProductDTO, error) {
	uc.log.WithContext(ctx).Infof("GetProduct: %v", id)
	return uc.repo.GetProduct(ctx, id)
}

// ListProducts lists all Products.
func (uc *InventoryUsecase) ListProducts(ctx context.Context, query *dto.ListProductsQuery) ([]*dto.ProductDTO, int32, error) {
	uc.log.WithContext(ctx).Infof("ListProducts: page=%d, page_size=%d", query.Page, query.PageSize)

	// Business logic: default pagination
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100 // Max page size
	}

	return uc.repo.ListProducts(ctx, query)
}

// UpdateProduct updates a Product.
func (uc *InventoryUsecase) UpdateProduct(ctx context.Context, req *dto.UpdateProductDTO) (*dto.ProductDTO, error) {
	uc.log.WithContext(ctx).Infof("UpdateProduct: %v", req.ID)
	
	// Business logic: validation
	if req.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	
	return uc.repo.UpdateProduct(ctx, req)
}

// DeleteProduct deletes a Product.
func (uc *InventoryUsecase) DeleteProduct(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteProduct: %v", id)
	return uc.repo.DeleteProduct(ctx, id)
}

// UpdateStock updates product stock.
func (uc *InventoryUsecase) UpdateStock(ctx context.Context, req *dto.UpdateStockDTO) (*dto.ProductDTO, error) {
	uc.log.WithContext(ctx).Infof("UpdateStock: id=%s, quantity=%d, operation=%s", req.ID, req.Quantity, req.Operation)
	
	// Business logic: validate operation
	if req.Operation != "add" && req.Operation != "subtract" {
		return nil, fmt.Errorf("invalid operation: %s", req.Operation)
	}
	
	return uc.repo.UpdateStock(ctx, req)
}


