package object_storage_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/poly-workshop/go-webmods/object_storage"
)

// Example demonstrates basic object storage usage with the local provider.
func Example() {
	storage, err := object_storage.NewObjectStorage(object_storage.Config{
		ProviderType: object_storage.ProviderLocal,
		ProviderConfig: object_storage.ProviderConfig{
			BasePath: "/tmp/storage",
		},
	})
	if err != nil {
		panic(err)
	}

	// Save a file
	content := bytes.NewReader([]byte("Hello, World!"))
	size, err := storage.Save("hello.txt", content)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Saved %d bytes\n", size)

	// Open and read the file
	obj, err := storage.Open("hello.txt")
	if err != nil {
		panic(err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Content: %s\n", string(data))
	// Output: Saved 13 bytes
	// Content: Hello, World!
}

// Example_minio demonstrates using MinIO or S3-compatible storage.
func Example_minio() {
	storage, err := object_storage.NewObjectStorage(object_storage.Config{
		ProviderType: object_storage.ProviderMinio,
		ProviderConfig: object_storage.ProviderConfig{
			Endpoint:  "minio.example.com:9000",
			Region:    "us-east-1",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			Bucket:    "mybucket",
			BasePath:  "uploads/",
		},
	})
	if err != nil {
		panic(err)
	}

	_ = storage
	fmt.Println("MinIO storage configured")
	// Output: MinIO storage configured
}

// Example_list demonstrates listing objects in a directory.
func Example_list() {
	storage, err := object_storage.NewObjectStorage(object_storage.Config{
		ProviderType: object_storage.ProviderLocal,
		ProviderConfig: object_storage.ProviderConfig{
			BasePath: "/tmp/storage",
		},
	})
	if err != nil {
		panic(err)
	}

	// Save some files
	storage.Save("photos/photo1.jpg", bytes.NewReader([]byte("photo1")))
	storage.Save("photos/photo2.jpg", bytes.NewReader([]byte("photo2")))

	// List files in directory
	objects, err := storage.List("photos/")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d files\n", len(objects))
	for _, obj := range objects {
		info, _ := obj.Stat()
		fmt.Printf("- %s (%d bytes)\n", info.Name(), info.Size())
	}
}

// Example_stat demonstrates getting file metadata without downloading.
func Example_stat() {
	storage, err := object_storage.NewObjectStorage(object_storage.Config{
		ProviderType: object_storage.ProviderLocal,
		ProviderConfig: object_storage.ProviderConfig{
			BasePath: "/tmp/storage",
		},
	})
	if err != nil {
		panic(err)
	}

	// Save a file
	storage.Save("document.pdf", bytes.NewReader([]byte("PDF content")))

	// Get file info
	info, err := storage.Stat("document.pdf")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", info.Name())
	fmt.Printf("Size: %d bytes\n", info.Size())
	fmt.Printf("Modified: %s\n", info.ModTime().Format("2006-01-02"))
}

// Example_delete demonstrates deleting a file from storage.
func Example_delete() {
	storage, err := object_storage.NewObjectStorage(object_storage.Config{
		ProviderType: object_storage.ProviderLocal,
		ProviderConfig: object_storage.ProviderConfig{
			BasePath: "/tmp/storage",
		},
	})
	if err != nil {
		panic(err)
	}

	// Save a file
	storage.Save("temp.txt", bytes.NewReader([]byte("temporary")))

	// Delete the file
	err = storage.Delete("temp.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("File deleted")
	// Output: File deleted
}
