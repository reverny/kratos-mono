# Layered Architecture - Inventory Service

## Overview

Service ถูกออกแบบด้วย Clean Architecture แบ่งเป็น 4 layers ชัดเจน:

```
Request/Response (Proto) 
    ↓
Service Layer (DTO conversion)
    ↓
Business Logic Layer (DTO)
    ↓
Data Layer (Entity ↔ DTO conversion)
```

## Layer Responsibilities

### 1. Service Layer (`internal/service/`)

**รับผิดชอบ:**
- รับ gRPC/HTTP requests (protobuf messages)
- แปลง Request → DTO
- เรียก Business Logic ผ่าน DTO
- แปลง DTO → Response
- **ไม่มี business logic**

**ตัวอย่าง:**
```go
func (s *InventoryService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
    // Convert request to DTO
    createDTO := &dto.CreateProductDTO{
        Name:        req.Name,
        Description: req.Description,
        SKU:         req.Sku,
        Price:       req.Price,
        Stock:       req.Stock,
    }

    // Call business logic with DTO
    productDTO, err := s.uc.CreateProduct(ctx, createDTO)
    if err != nil {
        return nil, err
    }

    // Convert DTO back to response
    return dtoToProto(productDTO), nil
}
```

### 2. DTO Layer (`internal/dto/`)

**รับผิดชอบ:**
- Define data structures สำหรับ Business Logic
- แยก business data ออกจาก transport layer (proto)
- ไม่มี dependencies กับ proto หรือ entity

**ตัวอย่าง:**
```go
type ProductDTO struct {
    ID          string
    Name        string
    Description string
    SKU         string
    Price       float64
    Stock       int32
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 3. Business Logic Layer (`internal/biz/`)

**รับผิดชอบ:**
- Business rules และ validation
- Orchestrate data operations
- ทำงานกับ DTO เท่านั้น
- **ไม่รู้จัก Proto หรือ Entity**

**ตัวอย่าง:**
```go
func (uc *InventoryUsecase) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*dto.ProductDTO, error) {
    // Business logic: validation
    if req.Stock < 0 {
        req.Stock = 0
    }
    
    // Call repository with DTO
    return uc.repo.CreateProduct(ctx, req)
}
```

### 4. Data Layer (`internal/data/`)

**รับผิดชอบ:**
- รับ DTO จาก Business Layer
- แปลง DTO → Entity
- ทำงานกับ database ผ่าน Entity
- แปลง Entity → DTO กลับไป
- **ไม่มี business logic**

**ตัวอย่าง:**
```go
func (r *inventoryRepo) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*dto.ProductDTO, error) {
    // Create entity from DTO
    productEntity := &entity.Product{
        ID:          uuid.New().String(),
        Name:        req.Name,
        Description: req.Description,
        SKU:         req.SKU,
        Price:       req.Price,
        Stock:       req.Stock,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    // Save entity to database
    // r.data.db.Create(productEntity)
    
    // Convert entity back to DTO
    return productEntity.ToDTO(), nil
}
```

### 5. Entity Layer (`internal/data/entity/`)

**รับผิดชอบ:**
- Define database structures
- Mapping กับ database tables
- Conversion methods (Entity ↔ DTO)

**ตัวอย่าง:**
```go
type Product struct {
    ID          string
    Name        string
    Description string
    SKU         string
    Price       float64
    Stock       int32
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (e *Product) ToDTO() *dto.ProductDTO {
    return &dto.ProductDTO{
        ID:          e.ID,
        Name:        e.Name,
        // ...
    }
}
```

## Data Flow

### Create Product Flow:

```
1. gRPC Request (CreateProductRequest) 
   ↓
2. Service Layer
   - Convert to CreateProductDTO
   ↓
3. Business Layer (DTO)
   - Validate stock >= 0
   - Pass DTO to Repository
   ↓
4. Data Layer
   - Convert DTO → Entity
   - Save Entity to DB
   - Convert Entity → DTO
   ↓
5. Business Layer (DTO)
   - Return DTO
   ↓
6. Service Layer
   - Convert DTO → Product (proto)
   ↓
7. gRPC Response (Product)
```

## Benefits

### ✅ Separation of Concerns
- แต่ละ layer มี responsibility ชัดเจน
- ง่ายต่อการ test แยกส่วน

### ✅ Independence
- Business logic ไม่ผูกกับ transport protocol (proto)
- Data layer ไม่ผูกกับ business rules
- เปลี่ยน database schema ไม่กระทบ business logic

### ✅ Flexibility
- แต่ง response data ได้ใน service layer
- Business logic ใช้ได้กับหลาย protocols (gRPC, REST, GraphQL)
- เปลี่ยน database ได้ง่าย (เปลี่ยนแค่ entity layer)

### ✅ Maintainability
- โค้ดอ่านง่าย เข้าใจได้ชัดเจน
- แก้ไขส่วนใดส่วนหนึ่งไม่กระทบส่วนอื่น
- เพิ่ม feature ใหม่ง่าย

## Best Practices

1. **Service Layer**
   - แปลง proto ↔ DTO เท่านั้น
   - ห้ามมี business logic
   - ห้ามเรียก repository โดยตรง

2. **Business Layer**
   - ทำงานกับ DTO เท่านั้น
   - ใส่ business rules ทั้งหมดที่นี่
   - ห้ามรู้จัก proto หรือ entity

3. **Data Layer**
   - แปลง DTO ↔ Entity
   - ทำงานกับ database ผ่าน entity
   - ห้ามมี business logic

4. **DTO & Entity**
   - DTO สำหรับ business logic
   - Entity สำหรับ database mapping
   - แยกกันชัดเจน ไม่ใช้ผสม

## Testing Strategy

### Unit Tests
- **Service**: Mock business layer, test conversion
- **Business**: Mock repository, test business rules
- **Data**: Mock database, test entity operations

### Integration Tests
- Test ทั้ง flow จาก service → database

## Migration from Old Code

ถ้ามี service อื่นที่ยังไม่ได้ปรับ สามารถทำทีละขั้นตอน:

1. สร้าง DTO และ Entity structs
2. เพิ่ม conversion ใน service layer
3. Update business layer ให้ใช้ DTO
4. Update data layer ให้แปลง DTO ↔ Entity
5. ทดสอบทุกขั้นตอน

## Example: Adding New Field

เมื่อต้องการเพิ่ม field ใหม่ (เช่น `category`):

1. **Proto** - เพิ่มใน .proto file
2. **DTO** - เพิ่ม `Category string` ใน ProductDTO
3. **Entity** - เพิ่ม `Category string` ใน Product entity
4. **Service** - Update conversion proto ↔ DTO
5. **Data** - Update conversion DTO ↔ Entity

โครงสร้างที่ชัดเจนทำให้รู้ว่าต้องแก้ที่ไหนบ้าง!
