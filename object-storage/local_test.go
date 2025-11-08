package object_storage

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestLocalObjectStorage(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create local storage instance
	storage, err := NewLocalObjectStorage(ProviderConfig{
		BasePath: tempDir,
	})
	if err != nil {
		t.Fatalf("Failed to create local storage: %v", err)
	}

	// Test Save
	testData := "Hello, World!"
	testPath := "test/file.txt"

	written, err := storage.Save(testPath, strings.NewReader(testData))
	if err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}
	if written != int64(len(testData)) {
		t.Fatalf("Expected %d bytes written, got %d", len(testData), written)
	}

	// Test Stat
	info, err := storage.Stat(testPath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() != int64(len(testData)) {
		t.Fatalf("Expected file size %d, got %d", len(testData), info.Size())
	}

	// Test Open and Read
	obj, err := storage.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() {
		if err := obj.Close(); err != nil {
			t.Logf("Failed to close object: %v", err)
		}
	}()

	buf := make([]byte, len(testData))
	n, err := obj.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read file: %v", err)
	}
	if n != len(testData) {
		t.Fatalf("Expected to read %d bytes, got %d", len(testData), n)
	}
	if string(buf) != testData {
		t.Fatalf("Expected data %q, got %q", testData, string(buf))
	}

	// Test Seek
	_, err = obj.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("Failed to seek: %v", err)
	}

	// Read again after seek
	buf2 := make([]byte, len(testData))
	_, err = obj.Read(buf2)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read after seek: %v", err)
	}
	if string(buf2) != testData {
		t.Fatalf("Expected data after seek %q, got %q", testData, string(buf2))
	}

	// Test List
	// Create another file for listing test
	_, err = storage.Save("test/file2.txt", strings.NewReader("test2"))
	if err != nil {
		t.Fatalf("Failed to save second file: %v", err)
	}

	objects, err := storage.List("test")
	if err != nil {
		t.Fatalf("Failed to list objects: %v", err)
	}
	if len(objects) != 2 {
		t.Fatalf("Expected 2 objects, got %d", len(objects))
	}

	// Clean up objects
	for _, obj := range objects {
		if err := obj.Close(); err != nil {
			t.Logf("Failed to close object: %v", err)
		}
	}

	// Test Delete
	err = storage.Delete(testPath)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	// Verify file is deleted
	_, err = storage.Stat(testPath)
	if !os.IsNotExist(err) {
		t.Fatalf("File should be deleted, but still exists")
	}
}

func TestLocalObjectStorageDefaultPath(t *testing.T) {
	// Test with default path (empty BasePath)
	storage, err := NewLocalObjectStorage(ProviderConfig{})
	if err != nil {
		t.Fatalf("Failed to create local storage with default path: %v", err)
	}

	// Clean up default data directory if it was created
	defer func() { _ = os.RemoveAll("./data") }()

	// Just test that we can save a file
	testData := "test"
	testPath := "default_test.txt"

	_, err = storage.Save(testPath, strings.NewReader(testData))
	if err != nil {
		t.Fatalf("Failed to save file with default path: %v", err)
	}

	// Clean up
	_ = storage.Delete(testPath)
}
