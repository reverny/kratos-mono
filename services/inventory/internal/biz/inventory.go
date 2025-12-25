package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/reverny/kratos-mono/gen/go/api/inventory/v1"
)

// InventoryRepo is a Inventory repo.
type InventoryRepo interface {
	CreateProduct(context.Context, *v1.CreateProductRequest) (*v1.Product, error)
	GetProduct(context.Context, string) (*v1.Product, error)
	ListProducts(context.Context, int32, int32) ([]*v1.Product, int32, error)
	UpdateProduct(context.Context, *v1.UpdateProductRequest) (*v1.Product, error)
	DeleteProduct(context.Context, string) error
	UpdateStock(context.Context, string, int32, string) (*v1.Product, error)
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
func (uc *InventoryUsecase) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
	uc.log.WithContext(ctx).Infof("CreateProduct: %v", req)
	return uc.repo.CreateProduct(ctx, req)
}

// GetProduct gets a Product by ID.
func (uc *InventoryUsecase) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
	uc.log.WithContext(ctx).Infof("GetProduct: %v", req.Id)
	return uc.repo.GetProduct(ctx, req.Id)
}

// ListProducts lists all Products.
func (uc *InventoryUsecase) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	uc.log.WithContext(ctx).Infof("ListProducts: page=%d, page_size=%d", req.Page, req.PageSize)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	products, total, err := uc.repo.ListProducts(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ListProductsResponse{
		Products: products,
		Total:    total,
	}, nil
}

// UpdateProduct updates a Product.
func (uc *InventoryUsecase) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Product, error) {
	uc.log.WithContext(ctx).Infof("UpdateProduct: %v", req)
	return uc.repo.UpdateProduct(ctx, req)
}

// DeleteProduct deletes a Product.
func (uc *InventoryUsecase) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*emptypb.Empty, error) {
	uc.log.WithContext(ctx).Infof("DeleteProduct: %v", req.Id)
	err := uc.repo.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateStock updates product stock.
func (uc *InventoryUsecase) UpdateStock(ctx context.Context, req *v1.UpdateStockRequest) (*v1.Product, error) {
	uc.log.WithContext(ctx).Infof("UpdateStock: id=%s, quantity=%d, operation=%s", req.Id, req.Quantity, req.Operation)
	return uc.repo.UpdateStock(ctx, req.Id, req.Quantity, req.Operation)
}
