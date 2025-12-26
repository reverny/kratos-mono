package data

import (
	"context"

	"github.com/reverny/kratos-mono/services/test/internal/data/entity"
	"github.com/reverny/kratos-mono/services/test/internal/dto"
)

func (r *testRepo) Create(ctx context.Context, req *dto.CreateTestDTO) (*dto.TestDTO, error) {
	// TODO: implement database create
	// Convert CreateDTO to entity
	ent := &entity.Test{
		Name: req.Name,
	}
	
	// Simulate DB insert
	ent.ID = 1
	
	return ent.ToDTO(), nil
}

func (r *testRepo) Get(ctx context.Context, id int64) (*dto.TestDTO, error) {
	// TODO: implement database get
	ent := &entity.Test{
		ID:   id,
		Name: "sample",
	}
	return ent.ToDTO(), nil
}

func (r *testRepo) List(ctx context.Context, query *dto.ListTestQuery) ([]*dto.TestDTO, int, error) {
	// TODO: implement database list with pagination
	entities := []*entity.Test{
		{ID: 1, Name: "sample1"},
		{ID: 2, Name: "sample2"},
	}
	
	dtos := make([]*dto.TestDTO, len(entities))
	for i, ent := range entities {
		dtos[i] = ent.ToDTO()
	}
	
	return dtos, len(dtos), nil
}

func (r *testRepo) Update(ctx context.Context, req *dto.UpdateTestDTO) (*dto.TestDTO, error) {
	// TODO: implement database update
	ent := &entity.Test{
		ID:   req.ID,
		Name: req.Name,
	}
	return ent.ToDTO(), nil
}

func (r *testRepo) Delete(ctx context.Context, id int64) error {
	// TODO: implement database delete
	return nil
}
