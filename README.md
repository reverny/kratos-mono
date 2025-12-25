# Kratos Mono-repo

โครงสร้าง mono-repo สำหรับ microservices ด้วย [Kratos framework](https://go-kratos.dev/)

## โครงสร้าง

```
kratos-mono/
├── api/                    # API definitions (proto files)
│   ├── common/            # Shared proto files
│   └── inventory/         # Inventory service API
├── gen/                   # Generated code
│   └── go/               # Generated Go code from proto
├── services/              # Microservices
│   └── inventory/        # Inventory service
├── pkg/                   # Shared packages
├── third_party/          # Third-party proto files
├── go.work               # Go workspace
├── go.mod                # Root go.mod
├── buf.yaml              # Buf configuration
└── Makefile              # Build commands
```

## เริ่มต้นใช้งาน

### ติดตั้ง Dependencies

```bash
# ติดตั้ง Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# ติดตั้ง Buf
go install github.com/bufbuild/buf/cmd/buf@latest

# ติดตั้ง protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
```

### Generate Proto Files

```bash
make api
```

### รัน Service

```bash
# รัน inventory service
cd services/inventory
kratos run
```

## คำสั่ง Make

- `make api` - Generate code จาก proto files
- `make build` - Build ทุก services
- `make test` - Run tests
- `make lint` - Run linters
- `make clean` - Clean generated files

## เพิ่ม Service ใหม่

```bash
# ใช้ Kratos CLI สร้าง service ใหม่
cd services
kratos new <service-name>

# เพิ่มใน go.work
# แล้วสร้าง API definition ใน api/<service-name>/
```
