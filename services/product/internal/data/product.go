package data

import (
	"context"

	"github.com/reverny/kratos-mono/services/product/internal/biz"
)

func (r *productRepo) Create(ctx context.Context, item *biz.Product) (*biz.Product, error) {
	// TODO: implement database create
	item.ID = 1
	return item, nil
}

func (r *productRepo) Get(ctx context.Context, id int64) (*biz.Product, error) {
	// TODO: implement database get
	return &biz.Product{
		ID:   id,
		Name: "sample",
	}, nil
}

func (r *productRepo) List(ctx context.Context, page, pageSize int) ([]*biz.Product, int, error) {
	// TODO: implement database list
	items := []*biz.Product{
		{ID: 1, Name: "sample1"},
		{ID: 2, Name: "sample2"},
	}
	return items, len(items), nil
}

func (r *productRepo) Update(ctx context.Context, item *biz.Product) (*biz.Product, error) {
	// TODO: implement database update
	return item, nil
}

func (r *productRepo) Delete(ctx context.Context, id int64) error {
	// TODO: implement database delete
	return nil
}
