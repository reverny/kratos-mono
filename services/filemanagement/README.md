# File Management Service

Microservice สำหรับจัดการการอัปโหลดไฟล์โดยใช้ Presigned URL Pattern

## Architecture

- **Pattern**: Presigned URL (Best Practice for file uploads)
- **Storage**: Local (Development), MinIO/S3 (Production - TODO)
- **Ports**: 
  - HTTP: 8005
  - gRPC: 9005

## API Endpoints

### 1. Request Upload URL
```
POST /api/v1/files/upload/request
```
Request presigned URL สำหรับอัปโหลดไฟล์

**Request:**
```json
{
  "file_name": "example.pdf",
  "content_type": "application/pdf",
  "file_size": 1024000,
  "description": "Example file"
}
```

**Response:**
```json
{
  "file_id": "a1b2c3d4e5f6...",
  "upload_url": "http://localhost:8005/files/upload/a1b2c3d4e5f6...",
  "download_url": "http://localhost:8005/files/a1b2c3d4e5f6.../example.pdf",
  "method": "PUT",
  "headers": {
    "Content-Type": "application/pdf"
  },
  "expires_in": 900
}
```

### 2. Confirm Upload (Optional)
```
POST /api/v1/files/upload/confirm
```
ยืนยันว่าไฟล์ถูกอัปโหลดสำเร็จ

**Request:**
```json
{
  "file_id": "a1b2c3d4e5f6..."
}
```

### 3. Get File Info
```
GET /api/v1/files/{file_id}
```
ดึงข้อมูลเมตาดาต้าของไฟล์

### 4. Delete File
```
DELETE /api/v1/files/{file_id}
```
ลบไฟล์

## Usage Example

### 1. Request Upload URL
```bash
curl -X POST http://localhost:8005/api/v1/files/upload/request \
  -H "Content-Type: application/json" \
  -d '{
    "file_name": "document.pdf",
    "content_type": "application/pdf",
    "file_size": 1024000,
    "description": "Important document"
  }'
```

### 2. Upload File to Presigned URL
```bash
curl -X PUT "<upload_url>" \
  -H "Content-Type: application/pdf" \
  --data-binary "@document.pdf"
```

### 3. Confirm Upload (Optional)
```bash
curl -X POST http://localhost:8005/api/v1/files/upload/confirm \
  -H "Content-Type: application/json" \
  -d '{"file_id": "<file_id>"}'
```

## Development

### Build
```bash
make build
```

### Run
```bash
make run
```

### Generate Wire
```bash
make wire
```

## Configuration

Edit `configs/config.yaml`:

```yaml
storage:
  type: local  # local, minio, s3
  base_url: http://localhost:8005/files
  upload_path: ./uploads
  # MinIO/S3 config (for production)
  # endpoint: localhost:9000
  # access_key: minioadmin
  # secret_key: minioadmin
  # bucket: files
  # use_ssl: false
```

## Integration with Other Services

Services อื่นๆ (test, product, user, inventory) สามารถเรียกใช้ file-management service ผ่าน gRPC:

```go
// Add file-management gRPC client
fileClient := filemanagementv1.NewFilemanagementClient(conn)

// Request upload URL
resp, err := fileClient.RequestUploadUrl(ctx, &filemanagementv1.RequestUploadUrlRequest{
    FileName: "image.jpg",
    ContentType: "image/jpeg",
    FileSize: 204800,
})
```

## TODO

- [ ] Implement MinIOStorage for production
- [ ] Implement S3Storage for AWS
- [ ] Add file metadata to database
- [ ] Add file size limits
- [ ] Add file type validation
- [ ] Add virus scanning integration
- [ ] Add CDN support
