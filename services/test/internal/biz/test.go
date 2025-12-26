package biz

import (
	"context"

	"github.com/reverny/kratos-mono/services/test/internal/dto"

	"github.com/go-kratos/kratos/v2/log"
)

type TestRepo interface {
	Create(context.Context, *dto.CreateTestDTO) (*dto.TestDTO, error)
	Get(context.Context, int64) (*dto.TestDTO, error)
	List(context.Context, *dto.ListTestQuery) ([]*dto.TestDTO, int, error)
	Update(context.Context, *dto.UpdateTestDTO) (*dto.TestDTO, error)
	Delete(context.Context, int64) error
}

type TestUseCase struct {
	repo TestRepo
	log  *log.Helper
}

func NewTestUseCase(repo TestRepo, logger log.Logger) *TestUseCase {
	return &TestUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *TestUseCase) Create(ctx context.Context, req *dto.CreateTestDTO) (*dto.TestDTO, error) {
	uc.log.WithContext(ctx).Infof("CreateTest: %v", req.Name)
	return uc.repo.Create(ctx, req)
}

func (uc *TestUseCase) Get(ctx context.Context, id int64) (*dto.TestDTO, error) {
	uc.log.WithContext(ctx).Infof("GetTest: %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *TestUseCase) List(ctx context.Context, query *dto.ListTestQuery) ([]*dto.TestDTO, int, error) {
	uc.log.WithContext(ctx).Infof("ListTest: page=%d, pageSize=%d", query.Page, query.PageSize)
	return uc.repo.List(ctx, query)
}

func (uc *TestUseCase) Update(ctx context.Context, req *dto.UpdateTestDTO) (*dto.TestDTO, error) {
	uc.log.WithContext(ctx).Infof("UpdateTest: %v", req)
	return uc.repo.Update(ctx, req)
}

func (uc *TestUseCase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteTest: %d", id)
	return uc.repo.Delete(ctx, id)
}
