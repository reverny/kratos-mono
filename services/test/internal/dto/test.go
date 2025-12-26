package dto

// TestDTO represents data transfer object for business logic layer
type TestDTO struct {
	ID   int64
	Name string
}

// CreateTestDTO for creating new test
type CreateTestDTO struct {
	Name string
}

// UpdateTestDTO for updating test
type UpdateTestDTO struct {
	ID   int64
	Name string
}

// ListTestQuery for list query parameters
type ListTestQuery struct {
	Page     int32
	PageSize int32
}
