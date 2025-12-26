package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	"github.com/reverny/kratos-mono/services/inventory/internal/biz"
	"github.com/reverny/kratos-mono/services/inventory/internal/data/entity"
	"github.com/reverny/kratos-mono/services/inventory/internal/dto"
)

type inventoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewInventoryRepo .
func NewInventoryRepo(data *Data, logger log.Logger) biz.InventoryRepo {
	return &inventoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *inventoryRepo) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*dto.ProductDTO, error) {
	now := time.Now()
	
	// Create entity from DTO
	productEntity := &entity.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// TODO: Save entity to database
	// r.data.db.Create(productEntity)
	
	r.log.Infof("Product created: %s", productEntity.ID)

	// Convert entity back to DTO
	return productEntity.ToDTO(), nil
}

func (r *inventoryRepo) GetProduct(ctx context.Context, id string) (*dto.ProductDTO, error) {
	// TODO: Get entity from database
	// var productEntity entity.Product
	// r.data.db.First(&productEntity, \"id = ?\", id)
	
	// Mock entity for now
	productEntity := &entity.Product{
		ID:          id,
		Name:        "Sample Product",
		Description: "This is a sample product",
		SKU:         "SKU-001",
		Price:       99.99,
		Stock:       100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Convert entity to DTO
	return productEntity.ToDTO(), nil
}

func (r *inventoryRepo) ListProducts(ctx context.Context, query *dto.ListProductsQuery) ([]*dto.ProductDTO, int32, error) {
	// TODO: Get entities from database with pagination
	// var entities []entity.Product
	// r.data.db.Offset(int((query.Page - 1) * query.PageSize)).Limit(int(query.PageSize)).Find(&entities)
	
	// Mock entities for now
	entities := []*entity.Product{
		{
			ID:          "1",
			Name:        "Product 1",
			Description: "Description 1",
			SKU:         "SKU-001",
			Price:       99.99,
			Stock:       100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Product 2",
			Description: "Description 2",
			SKU:         "SKU-002",
			Price:       149.99,
			Stock:       50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Convert entities to DTOs
	dtos := make([]*dto.ProductDTO, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}

	return dtos, int32(len(dtos)), nil
}

func (r *inventoryRepo) UpdateProduct(ctx context.Context, req *dto.UpdateProductDTO) (*dto.ProductDTO, error) {
	// TODO: Get existing entity from database
	// var productEntity entity.Product
	// r.data.db.First(&productEntity, \"id = ?\", req.ID)
	
	// Mock entity
	productEntity := &entity.Product{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UpdatedAt:   time.Now(),
	}

	// TODO: Update entity in database
	// r.data.db.Save(productEntity)

	r.log.Infof("Product updated: %s", productEntity.ID)

	// Convert entity to DTO
	return productEntity.ToDTO(), nil
}

func (r *inventoryRepo) DeleteProduct(ctx context.Context, id string) error {
	// TODO: Delete entity from database
	// r.data.db.Delete(&entity.Product{}, \"id = ?\", id)
	
	r.log.Infof("Product deleted: %s", id)
	return nil
}

func (r *inventoryRepo) UpdateStock(ctx context.Context, req *dto.UpdateStockDTO) (*dto.ProductDTO, error) {
	// TODO: Get entity from database
	// var productEntity entity.Product
	// r.data.db.First(&productEntity, \"id = ?\", req.ID)
	
	// Mock current stock
	currentStock := int32(100)
	
	var newStock int32
	switch req.Operation {
	case "add":
		newStock = currentStock + req.Quantity
	case "subtract":
		newStock = currentStock - req.Quantity
		if newStock < 0 {
			return nil, fmt.Errorf("insufficient stock")
		}
	default:
		return nil, fmt.Errorf("invalid operation: %s", req.Operation)
	}

	// Update entity
	productEntity := &entity.Product{
		ID:        req.ID,
		Stock:     newStock,
		UpdatedAt: time.Now(),
	}

	// TODO: Save entity to database
	// r.data.db.Save(productEntity)

	r.log.Infof("Stock updated for product %s: %d -> %d", req.ID, currentStock, newStock)

	// Convert entity to DTO
	return productEntity.ToDTO(), nil
}
