package entity

import (
	"github.com/reverny/kratos-mono/services/test/internal/dto"
)

// Test represents the database entity
type Test struct {
	ID   int64
	Name string
}

// ToDTO converts entity to DTO
func (e *Test) ToDTO() *dto.TestDTO {
	return &dto.TestDTO{
		ID:   e.ID,
		Name: e.Name,
	}
}

// FromDTO converts DTO to entity
func FromDTO(d *dto.TestDTO) *Test {
	return &Test{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromCreateDTO converts CreateDTO to entity
func FromCreateDTO(d *dto.CreateTestDTO) *Test {
	return &Test{
		Name: d.Name,
	}
}

// FromUpdateDTO converts UpdateDTO to entity
func FromUpdateDTO(d *dto.UpdateTestDTO) *Test {
	return &Test{
		ID:   d.ID,
		Name: d.Name,
	}
}
