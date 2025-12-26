package entity

import (
	"time"
	
	"github.com/reverny/kratos-mono/services/inventory/internal/dto"
)

// Product represents the database entity for product
type Product struct {
	ID          string
	Name        string
	Description string
	SKU         string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToDTO converts entity to DTO
func (e *Product) ToDTO() *dto.ProductDTO {
	return &dto.ProductDTO{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		SKU:         e.SKU,
		Price:       e.Price,
		Stock:       e.Stock,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// FromDTO converts DTO to entity
func FromDTO(d *dto.ProductDTO) *Product {
	return &Product{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		SKU:         d.SKU,
		Price:       d.Price,
		Stock:       d.Stock,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
