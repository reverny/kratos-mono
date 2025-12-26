# วิธีอัพเดท Swagger ที่แนะนำ

มี 2 วิธีที่ดีกว่าการใช้ Python script:

## วิธีที่ 1: ใช้ buf generate + แยก spec ด้วยตัวเอง (แนะนำ)

ปัญหาของ protoc-gen-openapi คือมันรวม spec ทุกอย่างไว้ในไฟล์เดียว

**แก้ไข:** ใช้ `grpc-gateway` แทน:

1. ติดตั้ง grpc-gateway plugin:
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

2. เพิ่มใน buf.gen.yaml:
```yaml
  - local: protoc-gen-openapiv2
    out: gen/openapi
    opt:
      - output_format=json
      - allow_merge=false
      - simple_operation_ids=true
```

3. รัน `buf generate` จะได้ไฟล์แยกสำหรับแต่ละ service

## วิธีที่ 2: ให้ Service generate OpenAPI spec ตอน runtime (ดีที่สุด)

ใช้ Kratos built-in OpenAPI generation:

1. ติดตั้ง swagger middleware:
```bash
go get github.com/go-kratos/swagger-api/openapiv2
```

2. เพิ่มใน http.go:
```go
import (
    "github.com/go-kratos/swagger-api/openapiv2"
)

func NewHTTPServer(...) *http.Server {
    // ...
    
    // Add OpenAPI handler
    openAPIHandler := openapiv2.NewHandler()
    srv.HandlePrefix("/q/", openAPIHandler)
    
    return srv
}
```

3. OpenAPI spec จะ available ที่ `/q/swagger.json`

## วิธีที่ 3: ใช้ Google Gnostic (ปัจจุบัน)

ถ้าต้องการใช้ต่อกับ protoc-gen-openapi:

ปัญหาคือ `protoc-gen-openapi` ไม่รองรับ google.api.http annotations ครบถ้วน

**วิธีแก้ชั่วคราว:** แก้ไข Python script ให้ parse proto file ได้ถูกต้อง

---

**คำแนะนำ:** ใช้วิธีที่ 2 เป็นวิธีที่ดีที่สุด เพราะ:
- ไม่ต้อง generate static files
- Spec จะถูกต้องเสมอ (ตาม runtime registration)
- ไม่ต้องพึ่ง Python หรือ script อื่นๆ
- น้อยขั้นตอน
