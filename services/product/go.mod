module github.com/reverny/kratos-mono/services/product

go 1.23.0

require (
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/google/wire v0.6.0
	github.com/reverny/kratos-mono v0.0.0
	go.uber.org/automaxprocs v1.6.0
	google.golang.org/protobuf v1.35.2
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241206012308-a4fef0638583 // indirect
	google.golang.org/grpc v1.69.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/reverny/kratos-mono => ../..
