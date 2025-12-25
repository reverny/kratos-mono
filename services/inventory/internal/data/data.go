package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/reverny/kratos-mono/services/inventory/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewInventoryRepo)

// Data .
type Data struct {
	// TODO: Add database connection, redis connection, etc.
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	data := &Data{
		log: log.NewHelper(logger),
	}

	return data, cleanup, nil
}
