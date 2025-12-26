// Example: File Upload with Presigned URL
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	pb "github.com/reverny/kratos-mono/gen/go/api/test/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTestClient(conn)

	// Upload file
	if err := uploadFileWithPresignedURL(client, "example.pdf"); err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
}

func uploadFileWithPresignedURL(client pb.TestClient, filePath string) error {
	ctx := context.Background()

	// Step 1: Open file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Step 2: Request presigned URL from API
	fmt.Println("📤 Requesting upload URL...")
	urlResp, err := client.RequestFileUpload(ctx, &pb.RequestFileUploadRequest{
		FileName:    fileInfo.Name(),
		ContentType: "application/pdf",
		FileSize:    fileInfo.Size(),
		Description: "Example file upload with presigned URL",
	})
	if err != nil {
		return fmt.Errorf("failed to request upload URL: %w", err)
	}

	fmt.Printf("✅ Got upload URL (expires in %d seconds)\n", urlResp.ExpiresIn)
	fmt.Printf("   File ID: %s\n", urlResp.FileId)
	fmt.Printf("   Upload URL: %s\n", urlResp.UploadUrl)

	// Step 3: Upload file directly to storage using presigned URL
	fmt.Println("\n📦 Uploading file to storage...")
	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req, err := http.NewRequest(urlResp.Method, urlResp.UploadUrl, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set headers from presigned URL response
	for key, value := range urlResp.Headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ File uploaded successfully to storage")

	// Step 4 (Optional): Confirm upload with API
	fmt.Println("\n✔️  Confirming upload...")
	confirmResp, err := client.ConfirmFileUpload(ctx, &pb.ConfirmFileUploadRequest{
		FileId: urlResp.FileId,
	})
	if err != nil {
		return fmt.Errorf("failed to confirm upload: %w", err)
	}

	fmt.Printf("✅ %s\n", confirmResp.Message)
	fmt.Printf("   Download URL: %s\n", confirmResp.FileUrl)

	// Step 5: Use file_id when creating/updating test item
	fmt.Println("\n📝 Creating test item with uploaded file...")
	testResp, err := client.CreateTest(ctx, &pb.CreateTestRequest{
		Name:   "Test with file",
		FileId: urlResp.FileId,
	})
	if err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}

	fmt.Println("\n🎉 Complete!")
	fmt.Printf("   Test ID: %d\n", testResp.Data.Id)
	fmt.Printf("   Name: %s\n", testResp.Data.Name)
	fmt.Printf("   File URL: %s\n", testResp.Data.FileUrl)

	return nil
}
