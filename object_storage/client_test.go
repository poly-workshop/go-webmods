package object_storage

import (
	"os"
	"strings"
	"testing"
)

func TestNewObjectStorageLocal(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_client_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Test creating local object storage through the factory function
	config := Config{
		ProviderType: ProviderLocal,
		ProviderConfig: ProviderConfig{
			BasePath: tempDir,
		},
	}

	storage, err := NewObjectStorage(config)
	if err != nil {
		t.Fatalf("Failed to create local object storage: %v", err)
	}

	// Test that it's actually a LocalObjectStorage
	localStorage, ok := storage.(*LocalObjectStorage)
	if !ok {
		t.Fatalf("Expected LocalObjectStorage, got %T", storage)
	}

	// Test basic functionality
	testData := "Hello from factory!"
	testPath := "factory_test.txt"

	written, err := localStorage.Save(testPath, strings.NewReader(testData))
	if err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}
	if written != int64(len(testData)) {
		t.Fatalf("Expected %d bytes written, got %d", len(testData), written)
	}

	// Clean up
	if err := localStorage.Delete(testPath); err != nil {
		t.Logf("Failed to delete test file: %v", err)
	}
}

func TestNewObjectStorageUnsupportedProvider(t *testing.T) {
	config := Config{
		ProviderType: "unsupported",
	}

	_, err := NewObjectStorage(config)
	if err == nil {
		t.Fatalf("Expected error for unsupported provider, got nil")
	}

	expectedError := "unsupported object storage provider: unsupported"
	if err.Error() != expectedError {
		t.Fatalf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestNewObjectStorageMinio(t *testing.T) {
	// Test creating MinIO object storage through the factory function
	config := Config{
		ProviderType: ProviderMinio,
		ProviderConfig: ProviderConfig{
			Endpoint:  "localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			Bucket:    "test-bucket",
			BasePath:  "test",
		},
	}

	storage, err := NewObjectStorage(config)
	if err != nil {
		t.Fatalf("Failed to create MinIO object storage: %v", err)
	}

	// Test that it's actually a MinioObjectStorage
	minioStorage, ok := storage.(*MinioObjectStorage)
	if !ok {
		t.Fatalf("Expected MinioObjectStorage, got %T", storage)
	}

	// Verify configuration was set correctly
	if minioStorage.bucket != "test-bucket" {
		t.Fatalf("Expected bucket 'test-bucket', got %q", minioStorage.bucket)
	}
	if minioStorage.basePath != "test" {
		t.Fatalf("Expected basePath 'test', got %q", minioStorage.basePath)
	}
}
