//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/reverny/kratos-mono/services/inventory/internal/biz"
	"github.com/reverny/kratos-mono/services/inventory/internal/conf"
	"github.com/reverny/kratos-mono/services/inventory/internal/data"
	"github.com/reverny/kratos-mono/services/inventory/internal/server"
	"github.com/reverny/kratos-mono/services/inventory/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
