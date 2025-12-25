package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	v1 "github.com/reverny/kratos-mono/gen/go/api/inventory/v1"
	"github.com/reverny/kratos-mono/services/inventory/internal/biz"
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

func (r *inventoryRepo) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
	now := time.Now().Format(time.RFC3339)
	product := &v1.Product{
		Id:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Sku:         req.Sku,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// TODO: Save to database
	r.log.Infof("Product created: %s", product.Id)

	return product, nil
}

func (r *inventoryRepo) GetProduct(ctx context.Context, id string) (*v1.Product, error) {
	// TODO: Get from database
	product := &v1.Product{
		Id:          id,
		Name:        "Sample Product",
		Description: "This is a sample product",
		Sku:         "SKU-001",
		Price:       99.99,
		Stock:       100,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	return product, nil
}

func (r *inventoryRepo) ListProducts(ctx context.Context, page, pageSize int32) ([]*v1.Product, int32, error) {
	// TODO: Get from database with pagination
	products := []*v1.Product{
		{
			Id:          "1",
			Name:        "Product 1",
			Description: "Description 1",
			Sku:         "SKU-001",
			Price:       99.99,
			Stock:       100,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			Id:          "2",
			Name:        "Product 2",
			Description: "Description 2",
			Sku:         "SKU-002",
			Price:       149.99,
			Stock:       50,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	return products, int32(len(products)), nil
}

func (r *inventoryRepo) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Product, error) {
	// TODO: Update in database
	product := &v1.Product{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	r.log.Infof("Product updated: %s", product.Id)

	return product, nil
}

func (r *inventoryRepo) DeleteProduct(ctx context.Context, id string) error {
	// TODO: Delete from database
	r.log.Infof("Product deleted: %s", id)
	return nil
}

func (r *inventoryRepo) UpdateStock(ctx context.Context, id string, quantity int32, operation string) (*v1.Product, error) {
	// TODO: Update stock in database

	var newStock int32
	currentStock := int32(100) // Get from database

	switch operation {
	case "add":
		newStock = currentStock + quantity
	case "subtract":
		newStock = currentStock - quantity
		if newStock < 0 {
			return nil, fmt.Errorf("insufficient stock")
		}
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	product := &v1.Product{
		Id:        id,
		Stock:     newStock,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	r.log.Infof("Stock updated for product %s: %d -> %d", id, currentStock, newStock)

	return product, nil
}
