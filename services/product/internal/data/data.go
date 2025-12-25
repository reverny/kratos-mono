package data

import (
	"github.com/reverny/kratos-mono/services/product/internal/biz"
	"github.com/reverny/kratos-mono/services/product/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewProductRepo)

type Data struct {
	log *log.Helper
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		log: log.NewHelper(logger),
	}, cleanup, nil
}

type productRepo struct {
	data *Data
	log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
