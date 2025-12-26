#!/bin/bash

set -e

# Check if service name is provided
if [ -z "$1" ]; then
    echo "Usage: make gen-service name=<service-name>"
    echo "Example: make gen-service name=product"
    exit 1
fi

SERVICE_NAME=$1
SERVICE_DIR="services/${SERVICE_NAME}"
API_DIR="api/${SERVICE_NAME}/v1"
PROTO_FILE="${API_DIR}/${SERVICE_NAME}.proto"
ROOT_DIR=$(pwd)

# Capitalize first letter function
capitalize() {
    echo "$1" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}'
}

SERVICE_NAME_UPPER=$(capitalize "$SERVICE_NAME")

# Check if service already exists
if [ -d "$SERVICE_DIR" ]; then
    echo "Error: Service '${SERVICE_NAME}' already exists"
    exit 1
fi

echo "Creating service: ${SERVICE_NAME}"

# Create API proto file
echo "Creating proto file..."
mkdir -p "$API_DIR"
cat > "$PROTO_FILE" <<EOF
syntax = "proto3";

package api.${SERVICE_NAME}.v1;

option go_package = "github.com/reverny/kratos-mono/gen/go/api/${SERVICE_NAME}/v1;v1";

import "google/api/annotations.proto";

service ${SERVICE_NAME_UPPER} {
  rpc Create${SERVICE_NAME_UPPER} (Create${SERVICE_NAME_UPPER}Request) returns (Create${SERVICE_NAME_UPPER}Reply) {
    option (google.api.http) = {
      post: "/api/v1/${SERVICE_NAME}"
      body: "*"
    };
  }
  rpc Get${SERVICE_NAME_UPPER} (Get${SERVICE_NAME_UPPER}Request) returns (Get${SERVICE_NAME_UPPER}Reply) {
    option (google.api.http) = {
      get: "/api/v1/${SERVICE_NAME}/{id}"
    };
  }
  rpc List${SERVICE_NAME_UPPER} (List${SERVICE_NAME_UPPER}Request) returns (List${SERVICE_NAME_UPPER}Reply) {
    option (google.api.http) = {
      get: "/api/v1/${SERVICE_NAME}"
    };
  }
  rpc Update${SERVICE_NAME_UPPER} (Update${SERVICE_NAME_UPPER}Request) returns (Update${SERVICE_NAME_UPPER}Reply) {
    option (google.api.http) = {
      put: "/api/v1/${SERVICE_NAME}/{id}"
      body: "*"
    };
  }
  rpc Delete${SERVICE_NAME_UPPER} (Delete${SERVICE_NAME_UPPER}Request) returns (Delete${SERVICE_NAME_UPPER}Reply) {
    option (google.api.http) = {
      delete: "/api/v1/${SERVICE_NAME}/{id}"
    };
  }
}

message ${SERVICE_NAME_UPPER}Item {
  int64 id = 1;
  string name = 2;
}

message Create${SERVICE_NAME_UPPER}Request {
  string name = 1;
}

message Create${SERVICE_NAME_UPPER}Reply {
  ${SERVICE_NAME_UPPER}Item data = 1;
}

message Get${SERVICE_NAME_UPPER}Request {
  int64 id = 1;
}

message Get${SERVICE_NAME_UPPER}Reply {
  ${SERVICE_NAME_UPPER}Item data = 1;
}

message List${SERVICE_NAME_UPPER}Request {
  int32 page = 1;
  int32 page_size = 2;
}

message List${SERVICE_NAME_UPPER}Reply {
  repeated ${SERVICE_NAME_UPPER}Item items = 1;
  int32 total = 2;
}

message Update${SERVICE_NAME_UPPER}Request {
  int64 id = 1;
  string name = 2;
}

message Update${SERVICE_NAME_UPPER}Reply {
  ${SERVICE_NAME_UPPER}Item data = 1;
}

message Delete${SERVICE_NAME_UPPER}Request {
  int64 id = 1;
}

message Delete${SERVICE_NAME_UPPER}Reply {
  bool success = 1;
}
EOF

# Generate proto code
echo "Generating proto code..."
buf generate api

# Create service directory structure
echo "Creating service directory structure..."
mkdir -p "${SERVICE_DIR}/cmd/${SERVICE_NAME}"
mkdir -p "${SERVICE_DIR}/configs"
mkdir -p "${SERVICE_DIR}/docs"
mkdir -p "${SERVICE_DIR}/internal/conf"
mkdir -p "${SERVICE_DIR}/internal/server"
mkdir -p "${SERVICE_DIR}/internal/service"
mkdir -p "${SERVICE_DIR}/internal/biz"
mkdir -p "${SERVICE_DIR}/internal/data"
mkdir -p "${SERVICE_DIR}/internal/data/entity"
mkdir -p "${SERVICE_DIR}/internal/dto"

# Get next available ports (increment by 1 from last service)
LAST_SERVICE=$(ls -d services/*/ 2>/dev/null | tail -1)
if [ -z "$LAST_SERVICE" ]; then
    HTTP_PORT=8000
    GRPC_PORT=9000
else
    LAST_CONFIG=$(find services/*/configs/config.yaml -type f 2>/dev/null | tail -1)
    if [ -n "$LAST_CONFIG" ]; then
        LAST_HTTP=$(grep -A 3 "http:" "$LAST_CONFIG" | grep "addr:" | sed 's/.*://g')
        HTTP_PORT=$((LAST_HTTP + 1))
        GRPC_PORT=$((HTTP_PORT + 1000))
    else
        HTTP_PORT=8000
        GRPC_PORT=9000
    fi
fi

# Create config.yaml
cat > "${SERVICE_DIR}/configs/config.yaml" <<EOF
server:
  http:
    addr: 0.0.0.0:${HTTP_PORT}
    timeout: 1s
  grpc:
    addr: 0.0.0.0:${GRPC_PORT}
    timeout: 1s
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/${SERVICE_NAME}?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
EOF

# Create conf.proto
cat > "${SERVICE_DIR}/internal/conf/conf.proto" <<EOF
syntax = "proto3";
package kratos.api;

option go_package = "github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf;conf";

import "google/protobuf/duration.proto";

message Bootstrap {
  Server server = 1;
  Data data = 2;
}

message Server {
  message HTTP {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  message GRPC {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  HTTP http = 1;
  GRPC grpc = 2;
}

message Data {
  message Database {
    string driver = 1;
    string source = 2;
  }
  message Redis {
    string addr = 1;
    google.protobuf.Duration read_timeout = 2;
    google.protobuf.Duration write_timeout = 3;
  }
  Database database = 1;
  Redis redis = 2;
}
EOF

# Generate conf proto
protoc --proto_path=. \
       --proto_path="${SERVICE_DIR}/third_party" \
       --go_out=paths=source_relative:"${SERVICE_DIR}/internal/conf" \
       "${SERVICE_DIR}/internal/conf/conf.proto" 2>/dev/null || \
protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=paths=source_relative:"${SERVICE_DIR}/internal/conf" \
       "${SERVICE_DIR}/internal/conf/conf.proto"

# Create main.go
cat > "${SERVICE_DIR}/cmd/${SERVICE_NAME}/main.go" <<EOF
package main

import (
	"flag"
	"os"

	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

var (
	Name    string = "${SERVICE_NAME}"
	Version string = "v1.0.0"

	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
EOF

# Create wire.go
cat > "${SERVICE_DIR}/cmd/${SERVICE_NAME}/wire.go" <<EOF
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/biz"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/data"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/server"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
EOF

# Create server.go
cat > "${SERVICE_DIR}/internal/server/server.go" <<EOF
package server

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
EOF

# Create grpc.go
cat > "${SERVICE_DIR}/internal/server/grpc.go" <<EOF
package server

import (
	v1 "github.com/reverny/kratos-mono/gen/go/api/${SERVICE_NAME}/v1"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(c *conf.Server, ${SERVICE_NAME}Svc *service.${SERVICE_NAME_UPPER}Service, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.Register${SERVICE_NAME_UPPER}Server(srv, ${SERVICE_NAME}Svc)
	return srv
}
EOF

# Create http.go
cat > "${SERVICE_DIR}/internal/server/http.go" <<EOF
package server

import (
	_ "embed"
	nethttp "net/http"

	v1 "github.com/reverny/kratos-mono/gen/go/api/${SERVICE_NAME}/v1"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

//go:embed swagger.html
var swaggerHTML []byte

func NewHTTPServer(c *conf.Server, ${SERVICE_NAME}Svc *service.${SERVICE_NAME_UPPER}Service, logger log.Logger) *http.Server {
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
	v1.Register${SERVICE_NAME_UPPER}HTTPServer(srv, ${SERVICE_NAME}Svc)
	
	// Serve Swagger UI
	srv.HandleFunc("/docs", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(swaggerHTML)
	})
	
	return srv
}
EOF

# Create service.go
cat > "${SERVICE_DIR}/internal/service/service.go" <<EOF
package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(New${SERVICE_NAME_UPPER}Service)
EOF

# Create service implementation
cat > "${SERVICE_DIR}/internal/service/${SERVICE_NAME}.go" <<EOF
package service

import (
	"context"

	pb "github.com/reverny/kratos-mono/gen/go/api/${SERVICE_NAME}/v1"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/biz"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/dto"
)

type ${SERVICE_NAME_UPPER}Service struct {
	pb.Unimplemented${SERVICE_NAME_UPPER}Server

	uc *biz.${SERVICE_NAME_UPPER}UseCase
}

func New${SERVICE_NAME_UPPER}Service(uc *biz.${SERVICE_NAME_UPPER}UseCase) *${SERVICE_NAME_UPPER}Service {
	return &${SERVICE_NAME_UPPER}Service{uc: uc}
}

func (s *${SERVICE_NAME_UPPER}Service) Create${SERVICE_NAME_UPPER}(ctx context.Context, req *pb.Create${SERVICE_NAME_UPPER}Request) (*pb.Create${SERVICE_NAME_UPPER}Reply, error) {
	item, err := s.uc.Create(ctx, &dto.Create${SERVICE_NAME_UPPER}DTO{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.Create${SERVICE_NAME_UPPER}Reply{
		Data: &pb.${SERVICE_NAME_UPPER}Item{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *${SERVICE_NAME_UPPER}Service) Get${SERVICE_NAME_UPPER}(ctx context.Context, req *pb.Get${SERVICE_NAME_UPPER}Request) (*pb.Get${SERVICE_NAME_UPPER}Reply, error) {
	item, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Get${SERVICE_NAME_UPPER}Reply{
		Data: &pb.${SERVICE_NAME_UPPER}Item{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *${SERVICE_NAME_UPPER}Service) List${SERVICE_NAME_UPPER}(ctx context.Context, req *pb.List${SERVICE_NAME_UPPER}Request) (*pb.List${SERVICE_NAME_UPPER}Reply, error) {
	items, total, err := s.uc.List(ctx, &dto.List${SERVICE_NAME_UPPER}Query{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	
	pbItems := make([]*pb.${SERVICE_NAME_UPPER}Item, len(items))
	for i, item := range items {
		pbItems[i] = &pb.${SERVICE_NAME_UPPER}Item{
			Id:   item.ID,
			Name: item.Name,
		}
	}
	
	return &pb.List${SERVICE_NAME_UPPER}Reply{
		Items: pbItems,
		Total: int32(total),
	}, nil
}

func (s *${SERVICE_NAME_UPPER}Service) Update${SERVICE_NAME_UPPER}(ctx context.Context, req *pb.Update${SERVICE_NAME_UPPER}Request) (*pb.Update${SERVICE_NAME_UPPER}Reply, error) {
	item, err := s.uc.Update(ctx, &dto.Update${SERVICE_NAME_UPPER}DTO{
		ID:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Update${SERVICE_NAME_UPPER}Reply{
		Data: &pb.${SERVICE_NAME_UPPER}Item{
			Id:   item.ID,
			Name: item.Name,
		},
	}, nil
}

func (s *${SERVICE_NAME_UPPER}Service) Delete${SERVICE_NAME_UPPER}(ctx context.Context, req *pb.Delete${SERVICE_NAME_UPPER}Request) (*pb.Delete${SERVICE_NAME_UPPER}Reply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Delete${SERVICE_NAME_UPPER}Reply{Success: true}, nil
}
EOF

# Create Swagger UI HTML
cat > "${SERVICE_DIR}/docs/swagger.html" <<'SWAGGEREOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SERVICENAME Service API</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        #swagger-ui { max-width: 1460px; margin: 0 auto; }
        .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                spec: {
                    openapi: '3.0.3',
                    info: {
                        title: 'SERVICENAME Service API',
                        version: '1.0.0',
                        description: 'SERVICENAME service API documentation'
                    },
                    servers: [
                        { url: window.location.origin, description: 'Current Server' }
                    ],
                    paths: {
                        '/api/v1/SERVICENAME': {
                            get: {
                                tags: ['SERVICENAME'],
                                summary: 'List items',
                                operationId: 'ListSERVICENAME',
                                parameters: [
                                    { name: 'page', in: 'query', schema: { type: 'integer', default: 1 } },
                                    { name: 'page_size', in: 'query', schema: { type: 'integer', default: 10 } }
                                ],
                                responses: {
                                    '200': {
                                        description: 'Success',
                                        content: {
                                            'application/json': {
                                                schema: { $ref: '#/components/schemas/ListSERVICENAMEReply' }
                                            }
                                        }
                                    }
                                }
                            },
                            post: {
                                tags: ['SERVICENAME'],
                                summary: 'Create item',
                                operationId: 'CreateSERVICENAME',
                                requestBody: {
                                    required: true,
                                    content: {
                                        'application/json': {
                                            schema: { $ref: '#/components/schemas/CreateSERVICENAMERequest' }
                                        }
                                    }
                                },
                                responses: {
                                    '200': {
                                        description: 'Success',
                                        content: {
                                            'application/json': {
                                                schema: { $ref: '#/components/schemas/CreateSERVICENAMEReply' }
                                            }
                                        }
                                    }
                                }
                            }
                        },
                        '/api/v1/SERVICENAME/{id}': {
                            get: {
                                tags: ['SERVICENAME'],
                                summary: 'Get item',
                                operationId: 'GetSERVICENAME',
                                parameters: [
                                    { name: 'id', in: 'path', required: true, schema: { type: 'integer', format: 'int64' } }
                                ],
                                responses: {
                                    '200': {
                                        description: 'Success',
                                        content: {
                                            'application/json': {
                                                schema: { $ref: '#/components/schemas/GetSERVICENAMEReply' }
                                            }
                                        }
                                    }
                                }
                            },
                            put: {
                                tags: ['SERVICENAME'],
                                summary: 'Update item',
                                operationId: 'UpdateSERVICENAME',
                                parameters: [
                                    { name: 'id', in: 'path', required: true, schema: { type: 'integer', format: 'int64' } }
                                ],
                                requestBody: {
                                    required: true,
                                    content: {
                                        'application/json': {
                                            schema: { $ref: '#/components/schemas/UpdateSERVICENAMERequest' }
                                        }
                                    }
                                },
                                responses: {
                                    '200': {
                                        description: 'Success',
                                        content: {
                                            'application/json': {
                                                schema: { $ref: '#/components/schemas/UpdateSERVICENAMEReply' }
                                            }
                                        }
                                    }
                                }
                            },
                            delete: {
                                tags: ['SERVICENAME'],
                                summary: 'Delete item',
                                operationId: 'DeleteSERVICENAME',
                                parameters: [
                                    { name: 'id', in: 'path', required: true, schema: { type: 'integer', format: 'int64' } }
                                ],
                                responses: {
                                    '200': {
                                        description: 'Success',
                                        content: {
                                            'application/json': {
                                                schema: { $ref: '#/components/schemas/DeleteSERVICENAMEReply' }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    components: {
                        schemas: {
                            SERVICENAMEItem: {
                                type: 'object',
                                properties: {
                                    id: { type: 'integer', format: 'int64' },
                                    name: { type: 'string' }
                                }
                            },
                            CreateSERVICENAMERequest: {
                                type: 'object',
                                required: ['name'],
                                properties: {
                                    name: { type: 'string' }
                                }
                            },
                            CreateSERVICENAMEReply: {
                                type: 'object',
                                properties: {
                                    data: { $ref: '#/components/schemas/SERVICENAMEItem' }
                                }
                            },
                            GetSERVICENAMEReply: {
                                type: 'object',
                                properties: {
                                    data: { $ref: '#/components/schemas/SERVICENAMEItem' }
                                }
                            },
                            ListSERVICENAMEReply: {
                                type: 'object',
                                properties: {
                                    items: {
                                        type: 'array',
                                        items: { $ref: '#/components/schemas/SERVICENAMEItem' }
                                    },
                                    total: { type: 'integer', format: 'int32' }
                                }
                            },
                            UpdateSERVICENAMERequest: {
                                type: 'object',
                                properties: {
                                    name: { type: 'string' }
                                }
                            },
                            UpdateSERVICENAMEReply: {
                                type: 'object',
                                properties: {
                                    data: { $ref: '#/components/schemas/SERVICENAMEItem' }
                                }
                            },
                            DeleteSERVICENAMEReply: {
                                type: 'object',
                                properties: {
                                    success: { type: 'boolean' }
                                }
                            }
                        },
                        securitySchemes: {
                            BearerAuth: {
                                type: 'http',
                                scheme: 'bearer',
                                bearerFormat: 'JWT',
                                description: 'Enter your Bearer token in the format: Bearer <token>'
                            }
                        }
                    },
                    security: [
                        {
                            BearerAuth: []
                        }
                    ]
                },
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
SWAGGEREOF

# Replace SERVICENAME with actual service name (uppercase first letter)
sed -i.bak "s/SERVICENAME/${SERVICE_NAME_UPPER}/g" "${SERVICE_DIR}/docs/swagger.html"
rm "${SERVICE_DIR}/docs/swagger.html.bak"

# Copy swagger.html to internal/server for embedding
cp "${SERVICE_DIR}/docs/swagger.html" "${SERVICE_DIR}/internal/server/swagger.html"

# Create main.go
cat > "${SERVICE_DIR}/internal/biz/biz.go" <<EOF
package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(New${SERVICE_NAME_UPPER}UseCase)
EOF

# Create biz implementation
cat > "${SERVICE_DIR}/internal/biz/${SERVICE_NAME}.go" <<EOF
package biz

import (
	"context"

	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/dto"

	"github.com/go-kratos/kratos/v2/log"
)

type ${SERVICE_NAME_UPPER}Repo interface {
	Create(context.Context, *dto.Create${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error)
	Get(context.Context, int64) (*dto.${SERVICE_NAME_UPPER}DTO, error)
	List(context.Context, *dto.List${SERVICE_NAME_UPPER}Query) ([]*dto.${SERVICE_NAME_UPPER}DTO, int, error)
	Update(context.Context, *dto.Update${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error)
	Delete(context.Context, int64) error
}

type ${SERVICE_NAME_UPPER}UseCase struct {
	repo ${SERVICE_NAME_UPPER}Repo
	log  *log.Helper
}

func New${SERVICE_NAME_UPPER}UseCase(repo ${SERVICE_NAME_UPPER}Repo, logger log.Logger) *${SERVICE_NAME_UPPER}UseCase {
	return &${SERVICE_NAME_UPPER}UseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *${SERVICE_NAME_UPPER}UseCase) Create(ctx context.Context, req *dto.Create${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	uc.log.WithContext(ctx).Infof("Create${SERVICE_NAME_UPPER}: %v", req.Name)
	return uc.repo.Create(ctx, req)
}

func (uc *${SERVICE_NAME_UPPER}UseCase) Get(ctx context.Context, id int64) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	uc.log.WithContext(ctx).Infof("Get${SERVICE_NAME_UPPER}: %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *${SERVICE_NAME_UPPER}UseCase) List(ctx context.Context, query *dto.List${SERVICE_NAME_UPPER}Query) ([]*dto.${SERVICE_NAME_UPPER}DTO, int, error) {
	uc.log.WithContext(ctx).Infof("List${SERVICE_NAME_UPPER}: page=%d, pageSize=%d", query.Page, query.PageSize)
	return uc.repo.List(ctx, query)
}

func (uc *${SERVICE_NAME_UPPER}UseCase) Update(ctx context.Context, req *dto.Update${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	uc.log.WithContext(ctx).Infof("Update${SERVICE_NAME_UPPER}: %v", req)
	return uc.repo.Update(ctx, req)
}

func (uc *${SERVICE_NAME_UPPER}UseCase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("Delete${SERVICE_NAME_UPPER}: %d", id)
	return uc.repo.Delete(ctx, id)
}
EOF

# Create data.go
cat > "${SERVICE_DIR}/internal/data/data.go" <<EOF
package data

import (
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/biz"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, New${SERVICE_NAME_UPPER}Repo)

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

type ${SERVICE_NAME}Repo struct {
	data *Data
	log  *log.Helper
}

func New${SERVICE_NAME_UPPER}Repo(data *Data, logger log.Logger) biz.${SERVICE_NAME_UPPER}Repo {
	return &${SERVICE_NAME}Repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
EOF

# Create data implementation
cat > "${SERVICE_DIR}/internal/data/${SERVICE_NAME}.go" <<EOF
package data

import (
	"context"

	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/data/entity"
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/dto"
)

func (r *${SERVICE_NAME}Repo) Create(ctx context.Context, req *dto.Create${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	// TODO: implement database create
	// Convert CreateDTO to entity
	ent := &entity.${SERVICE_NAME_UPPER}{
		Name: req.Name,
	}
	
	// Simulate DB insert
	ent.ID = 1
	
	return ent.ToDTO(), nil
}

func (r *${SERVICE_NAME}Repo) Get(ctx context.Context, id int64) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	// TODO: implement database get
	ent := &entity.${SERVICE_NAME_UPPER}{
		ID:   id,
		Name: "sample",
	}
	return ent.ToDTO(), nil
}

func (r *${SERVICE_NAME}Repo) List(ctx context.Context, query *dto.List${SERVICE_NAME_UPPER}Query) ([]*dto.${SERVICE_NAME_UPPER}DTO, int, error) {
	// TODO: implement database list with pagination
	entities := []*entity.${SERVICE_NAME_UPPER}{
		{ID: 1, Name: "sample1"},
		{ID: 2, Name: "sample2"},
	}
	
	dtos := make([]*dto.${SERVICE_NAME_UPPER}DTO, len(entities))
	for i, ent := range entities {
		dtos[i] = ent.ToDTO()
	}
	
	return dtos, len(dtos), nil
}

func (r *${SERVICE_NAME}Repo) Update(ctx context.Context, req *dto.Update${SERVICE_NAME_UPPER}DTO) (*dto.${SERVICE_NAME_UPPER}DTO, error) {
	// TODO: implement database update
	ent := &entity.${SERVICE_NAME_UPPER}{
		ID:   req.ID,
		Name: req.Name,
	}
	return ent.ToDTO(), nil
}

func (r *${SERVICE_NAME}Repo) Delete(ctx context.Context, id int64) error {
	// TODO: implement database delete
	return nil
}
EOF

# Create DTO files
echo "Creating DTO layer..."
cat > "${SERVICE_DIR}/internal/dto/${SERVICE_NAME}.go" <<EOF
package dto

// ${SERVICE_NAME_UPPER}DTO represents data transfer object for business logic layer
type ${SERVICE_NAME_UPPER}DTO struct {
	ID   int64
	Name string
}

// Create${SERVICE_NAME_UPPER}DTO for creating new ${SERVICE_NAME}
type Create${SERVICE_NAME_UPPER}DTO struct {
	Name string
}

// Update${SERVICE_NAME_UPPER}DTO for updating ${SERVICE_NAME}
type Update${SERVICE_NAME_UPPER}DTO struct {
	ID   int64
	Name string
}

// List${SERVICE_NAME_UPPER}Query for list query parameters
type List${SERVICE_NAME_UPPER}Query struct {
	Page     int32
	PageSize int32
}
EOF

# Create Entity files
echo "Creating Entity layer..."
cat > "${SERVICE_DIR}/internal/data/entity/${SERVICE_NAME}.go" <<EOF
package entity

import (
	"github.com/reverny/kratos-mono/services/${SERVICE_NAME}/internal/dto"
)

// ${SERVICE_NAME_UPPER} represents the database entity
type ${SERVICE_NAME_UPPER} struct {
	ID   int64
	Name string
}

// ToDTO converts entity to DTO
func (e *${SERVICE_NAME_UPPER}) ToDTO() *dto.${SERVICE_NAME_UPPER}DTO {
	return &dto.${SERVICE_NAME_UPPER}DTO{
		ID:   e.ID,
		Name: e.Name,
	}
}

// FromDTO converts DTO to entity
func FromDTO(d *dto.${SERVICE_NAME_UPPER}DTO) *${SERVICE_NAME_UPPER} {
	return &${SERVICE_NAME_UPPER}{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromCreateDTO converts CreateDTO to entity
func FromCreateDTO(d *dto.Create${SERVICE_NAME_UPPER}DTO) *${SERVICE_NAME_UPPER} {
	return &${SERVICE_NAME_UPPER}{
		Name: d.Name,
	}
}

// FromUpdateDTO converts UpdateDTO to entity
func FromUpdateDTO(d *dto.Update${SERVICE_NAME_UPPER}DTO) *${SERVICE_NAME_UPPER} {
	return &${SERVICE_NAME_UPPER}{
		ID:   d.ID,
		Name: d.Name,
	}
}
EOF

# Create go.mod
cat > "${SERVICE_DIR}/go.mod" <<EOF
module github.com/reverny/kratos-mono/services/${SERVICE_NAME}

go 1.23.0

require (
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/google/wire v0.6.0
	github.com/reverny/kratos-mono v0.0.0
	go.uber.org/automaxprocs v1.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.35.2
)

replace github.com/reverny/kratos-mono => ../..
EOF

# Generate conf proto first
echo "Generating conf proto..."
cd "${SERVICE_DIR}" && \
protoc --proto_path=. \
       --proto_path=../../third_party \
       --go_out=. \
       --go_opt=paths=source_relative \
       internal/conf/conf.proto && \
cd "${ROOT_DIR}"

# Add to go.work if exists
if [ -f "go.work" ]; then
    if ! grep -q "./services/${SERVICE_NAME}" go.work; then
        echo "Adding to go.work..."
        # Use go work use to add the module
        go work use "./services/${SERVICE_NAME}"
    fi
fi

# Run go work sync first
echo "Syncing workspace..."
go work sync

# Run go mod tidy
echo "Running go mod tidy..."
cd "${SERVICE_DIR}" && go mod tidy

# Run wire
echo "Running wire..."
cd "cmd/${SERVICE_NAME}" && wire

echo ""
echo "✅ Service '${SERVICE_NAME}' created successfully!"
echo ""
echo "Service details:"
echo "  - HTTP Port: ${HTTP_PORT}"
echo "  - gRPC Port: ${GRPC_PORT}"
echo "  - Proto: api/${SERVICE_NAME}/v1/${SERVICE_NAME}.proto"
echo "  - Location: services/${SERVICE_NAME}"
echo ""
echo "Next steps:"
echo "  1. Build: make build"
echo "  2. Run:   make run-${SERVICE_NAME}"
echo ""
