package dto

import "time"

// ProductDTO represents product data transfer object for business logic layer
type ProductDTO struct {
	ID          string
	Name        string
	Description string
	SKU         string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateProductDTO for creating new product
type CreateProductDTO struct {
	Name        string
	Description string
	SKU         string
	Price       float64
	Stock       int32
}

// UpdateProductDTO for updating product
type UpdateProductDTO struct {
	ID          string
	Name        string
	Description string
	Price       float64
}

// UpdateStockDTO for stock operations
type UpdateStockDTO struct {
	ID        string
	Quantity  int32
	Operation string // "add" or "subtract"
}

// ListProductsQuery for list query parameters
type ListProductsQuery struct {
	Page     int32
	PageSize int32
}
