package server

import (
	_ "embed"
	nethttp "net/http"

	v1 "github.com/reverny/kratos-mono/gen/go/api/product/v1"
	"github.com/reverny/kratos-mono/services/product/internal/conf"
	"github.com/reverny/kratos-mono/services/product/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

//go:embed swagger.html
var swaggerHTML []byte

func NewHTTPServer(c *conf.Server, productSvc *service.ProductService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterProductHTTPServer(srv, productSvc)
	
	// Serve Swagger UI
	srv.HandleFunc("/docs", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(swaggerHTML)
	})
	
	return srv
}
