# File Upload Architecture - Presigned URL Pattern

## 🏗️ Overview

เปลี่ยนจาก gRPC Streaming/Base64 JSON เป็น **Presigned URL Pattern** ซึ่งเป็น best practice สำหรับ file upload ใน production

## 📋 API Workflow

### 1. Request Upload URL
```http
POST /api/v1/test/files/upload-url
Content-Type: application/json

{
  "file_name": "document.pdf",
  "content_type": "application/pdf",
  "file_size": 1024000,
  "description": "Optional description"
}
```

**Response:**
```json
{
  "file_id": "abc123...",
  "upload_url": "https://storage.example.com/upload/abc123...",
  "download_url": "https://storage.example.com/files/abc123...",
  "expires_in": 900,
  "method": "PUT",
  "headers": {
    "Content-Type": "application/pdf"
  }
}
```

### 2. Upload File to Storage (Client → Storage)
```http
PUT https://storage.example.com/upload/abc123...
Content-Type: application/pdf

<binary file data>
```

### 3. Confirm Upload (Optional)
```http
POST /api/v1/test/files/confirm
Content-Type: application/json

{
  "file_id": "abc123..."
}
```

### 4. Use file_id in Business API
```http
POST /api/v1/test
Content-Type: application/json

{
  "name": "Test item",
  "file_id": "abc123..."
}
```

## 🎯 Benefits

✅ **Scalable**: File upload ไม่ผ่าน API server  
✅ **Fast**: Client upload ตรงไป Storage (S3/MinIO)  
✅ **Bandwidth Efficient**: ไม่กิน bandwidth ของ API  
✅ **Standard Pattern**: ใช้ในระบบใหญ่ทั้งหมด (AWS, GCP, Azure)  
✅ **Storage Agnostic**: เปลี่ยน storage backend ได้ง่าย  

## 🔧 Storage Implementations

### Current: LocalFileStorage (Development)
- ใช้สำหรับ development/testing
- เก็บไฟล์ local file system

### Production: MinIOStorage (Recommended)
```go
// TODO: Implement MinIOStorage
type MinIOStorage struct {
    client *minio.Client
    bucket string
}
```

### Cloud: S3Storage
```go
// TODO: Implement S3Storage
type S3Storage struct {
    client *s3.Client
    bucket string
}
```

## 📂 File Structure

```
services/test/
├── internal/
│   ├── biz/
│   │   ├── file_upload.go          # File upload use case
│   │   └── test.go
│   ├── data/
│   │   ├── file_storage.go         # Storage interface & implementations
│   │   └── data.go
│   └── service/
│       └── test.go                 # gRPC service with file endpoints
└── examples/
    └── presigned_upload_example.go # Example client code
```

## 🚀 Next Steps

1. [ ] Implement MinIOStorage for production use
2. [ ] Add file metadata to database
3. [ ] Implement file cleanup/expiry
4. [ ] Add file size limits
5. [ ] Add virus scanning
6. [ ] Add image processing (resize, thumbnail)
