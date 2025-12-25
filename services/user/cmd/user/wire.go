//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/reverny/kratos-mono/services/user/internal/biz"
	"github.com/reverny/kratos-mono/services/user/internal/conf"
	"github.com/reverny/kratos-mono/services/user/internal/data"
	"github.com/reverny/kratos-mono/services/user/internal/server"
	"github.com/reverny/kratos-mono/services/user/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
