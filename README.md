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
- `make update-swagger` - อัพเดท Swagger documentation จาก proto files

## อัพเดท Swagger Documentation

เมื่อมีการแก้ไขหรืออัพเดท proto files ต้องทำตามขั้นตอนดังนี้:

### 1. Generate API Code และ OpenAPI Specs

```bash
buf generate
```

หรือ

```bash
make api
```

คำสั่งนี้จะ:
- Generate Go code จาก proto files
- Generate OpenAPI specs (`.swagger.json`) สำหรับแต่ละ service โดยใช้ `grpc-gateway`
- ไฟล์ OpenAPI จะอยู่ที่ `gen/openapi/{service}/v1/{service}.swagger.json`

### 2. อัพเดท Swagger HTML Files

```bash
make update-swagger
```

คำสั่งนี้จะ:
- รัน `buf generate` อัตโนมัติ
- อ่าน OpenAPI specs ที่ generate แล้วจาก `gen/openapi/`
- อัพเดท swagger.html ใน `services/{service}/internal/server/swagger.html`
- อัพเดท swagger.html ใน `services/{service}/docs/swagger.html`

### ขั้นตอนย่อ

```bash
# หลังจากแก้ไข proto files
make update-swagger   # Generate code และ update swagger docs
make build           # Build services เพื่อ embed swagger HTML
```

### หมายเหตุ

- ใช้ **grpc-gateway** (`protoc-gen-openapiv2`) สำหรับ generate OpenAPI specs
- ไม่ต้องพึ่ง Python script ในการ parse proto files
- OpenAPI specs ถูก generate โดยตรงจาก proto annotations
- ทุก endpoint จะครบถ้วนตามที่กำหนดใน proto files

## เพิ่ม Service ใหม่

### ขั้นตอนการเพิ่ม Service

```bash
# 1. สร้าง service ด้วย Kratos CLI
cd services
kratos new <service-name>

# 2. เพิ่มใน go.work
cd ..
go work use ./services/<service-name>

# 3. สร้าง API definition (proto file)
mkdir -p api/<service-name>/v1
# สร้างไฟล์ api/<service-name>/v1/<service-name>.proto
```

### ตัวอย่าง Proto File

```protobuf
syntax = "proto3";

package api.<service-name>.v1;

option go_package = "github.com/reverny/kratos-mono/gen/go/api/<service-name>/v1;v1";

import "google/api/annotations.proto";

service <ServiceName> {
  rpc Create<Entity> (Create<Entity>Request) returns (Create<Entity>Reply) {
    option (google.api.http) = {
      post: "/api/v1/<entity>"
      body: "*"
    };
  }
  
  rpc Get<Entity> (Get<Entity>Request) returns (Get<Entity>Reply) {
    option (google.api.http) = {
      get: "/api/v1/<entity>/{id}"
    };
  }
}

// Define your messages here...
```

### Generate Code และ Swagger

```bash
# 4. Generate API code และ swagger
make update-swagger

# 5. Build service
make build
```

### หมายเหตุ

- ✅ **Auto-detection**: System จะ detect services ใหม่อัตโนมัติจาก folder `services/` และ `gen/openapi/`
- ✅ **No configuration needed**: ไม่ต้องแก้ไข Makefile หรือ script อื่นๆ
- ✅ **Swagger auto-generated**: OpenAPI specs จะถูก generate และ update อัตโนมัติ
- 📝 อย่าลืมเพิ่ม service ใหม่ใน `go.work` และ update dependencies
