//go:build wireinject
// +build wireinject

package main

import (
	"github.com/reverny/kratos-mono/services/product/internal/biz"
	"github.com/reverny/kratos-mono/services/product/internal/conf"
	"github.com/reverny/kratos-mono/services/product/internal/data"
	"github.com/reverny/kratos-mono/services/product/internal/server"
	"github.com/reverny/kratos-mono/services/product/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
