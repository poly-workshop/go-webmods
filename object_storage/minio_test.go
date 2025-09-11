package object_storage

import (
	"testing"
)

// TestMinioObjectStorageUnit tests MinIO object storage functionality without requiring a live MinIO server
// Note: These are unit tests that test the basic structure and error handling
func TestMinioObjectStorageUnit(t *testing.T) {
	// Test creating MinIO storage instance
	storage, err := NewMinioObjectStorage(ProviderConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "test-bucket",
		BasePath:  "test",
	})
	if err != nil {
		t.Fatalf("Failed to create MinIO storage: %v", err)
	}

	// Test getObjectPath method
	testPath := "file.txt"
	expectedPath := "test/file.txt"
	actualPath := storage.getObjectPath(testPath)
	if actualPath != expectedPath {
		t.Fatalf("Expected object path %q, got %q", expectedPath, actualPath)
	}

	// Test getObjectPath with empty basePath
	storageNoBase, err := NewMinioObjectStorage(ProviderConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "test-bucket",
		BasePath:  "",
	})
	if err != nil {
		t.Fatalf("Failed to create MinIO storage with empty basePath: %v", err)
	}

	actualPathNoBase := storageNoBase.getObjectPath(testPath)
	if actualPathNoBase != testPath {
		t.Fatalf("Expected object path %q with empty basePath, got %q", testPath, actualPathNoBase)
	}
}

// TestMinioObjectUnit tests MinIO object functionality
func TestMinioObjectUnit(t *testing.T) {
	// Create a MinIO object (without live connection)
	obj := &MinioObject{
		bucket:     "test-bucket",
		key:        "test/file.txt",
		rangeStart: 0,
	}

	// Test Close method
	err := obj.Close()
	if err != nil {
		t.Fatalf("Failed to close MinIO object: %v", err)
	}

	// Verify client is nil after close
	if obj.client != nil {
		t.Fatalf("Expected client to be nil after close")
	}
}

// Note: Integration tests would require a live MinIO server, which is not available in this test environment.
// In a real-world scenario, you would add integration tests similar to TestLocalObjectStorage that:
// 1. Set up a test MinIO server (using testcontainers or similar)
// 2. Test all ObjectStorage interface methods (Save, List, Open, Stat, Delete)
// 3. Test Object interface methods (Read, Seek, Close, Stat)
// 4. Test edge cases and error conditions
//
// Example integration test structure:
//
// func TestMinioObjectStorageIntegration(t *testing.T) {
//     if testing.Short() {
//         t.Skip("Skipping integration test in short mode")
//     }
//
//     // Setup test MinIO server
//     storage := setupTestMinioServer(t)
//     defer cleanupTestMinioServer(t)
//
//     // Test Save
//     testData := "Hello, MinIO!"
//     testPath := "test/file.txt"
//
//     written, err := storage.Save(testPath, strings.NewReader(testData))
//     if err != nil {
//         t.Fatalf("Failed to save file: %v", err)
//     }
//     if written != int64(len(testData)) {
//         t.Fatalf("Expected %d bytes written, got %d", len(testData), written)
//     }
//
//     // Test Stat
//     info, err := storage.Stat(testPath)
//     if err != nil {
//         t.Fatalf("Failed to stat file: %v", err)
//     }
//     if info.Size() != int64(len(testData)) {
//         t.Fatalf("Expected file size %d, got %d", len(testData), info.Size())
//     }
//
//     // Test Open and Read
//     obj, err := storage.Open(testPath)
//     if err != nil {
//         t.Fatalf("Failed to open file: %v", err)
//     }
//     defer obj.Close()
//
//     buf := make([]byte, len(testData))
//     n, err := obj.Read(buf)
//     if err != nil && err != io.EOF {
//         t.Fatalf("Failed to read file: %v", err)
//     }
//     if n != len(testData) {
//         t.Fatalf("Expected to read %d bytes, got %d", len(testData), n)
//     }
//     if string(buf) != testData {
//         t.Fatalf("Expected data %q, got %q", testData, string(buf))
//     }
//
//     // Test Seek
//     _, err = obj.Seek(0, io.SeekStart)
//     if err != nil {
//         t.Fatalf("Failed to seek: %v", err)
//     }
//
//     // Test List
//     objects, err := storage.List("test")
//     if err != nil {
//         t.Fatalf("Failed to list objects: %v", err)
//     }
//     if len(objects) != 1 {
//         t.Fatalf("Expected 1 object, got %d", len(objects))
//     }
//
//     // Test Delete
//     err = storage.Delete(testPath)
//     if err != nil {
//         t.Fatalf("Failed to delete file: %v", err)
//     }
// }

// TestNewMinioObjectStorageInvalidConfig tests error handling for invalid configuration
func TestNewMinioObjectStorageInvalidConfig(t *testing.T) {
	// Test with invalid endpoint (empty)
	_, err := NewMinioObjectStorage(ProviderConfig{
		Endpoint:  "",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "test-bucket",
	})
	if err == nil {
		t.Fatalf("Expected error for empty endpoint, got nil")
	}
}
