package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Product struct {
	ID   int64
	Name string
}

type ProductRepo interface {
	Create(context.Context, *Product) (*Product, error)
	Get(context.Context, int64) (*Product, error)
	List(context.Context, int, int) ([]*Product, int, error)
	Update(context.Context, *Product) (*Product, error)
	Delete(context.Context, int64) error
}

type ProductUseCase struct {
	repo ProductRepo
	log  *log.Helper
}

func NewProductUseCase(repo ProductRepo, logger log.Logger) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *ProductUseCase) Create(ctx context.Context, item *Product) (*Product, error) {
	uc.log.WithContext(ctx).Infof("CreateProduct: %v", item.Name)
	return uc.repo.Create(ctx, item)
}

func (uc *ProductUseCase) Get(ctx context.Context, id int64) (*Product, error) {
	uc.log.WithContext(ctx).Infof("GetProduct: %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *ProductUseCase) List(ctx context.Context, page, pageSize int) ([]*Product, int, error) {
	uc.log.WithContext(ctx).Infof("ListProduct: page=%d, pageSize=%d", page, pageSize)
	return uc.repo.List(ctx, page, pageSize)
}

func (uc *ProductUseCase) Update(ctx context.Context, item *Product) (*Product, error) {
	uc.log.WithContext(ctx).Infof("UpdateProduct: %v", item)
	return uc.repo.Update(ctx, item)
}

func (uc *ProductUseCase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteProduct: %d", id)
	return uc.repo.Delete(ctx, id)
}
