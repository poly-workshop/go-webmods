// Package object_storage provides a unified interface for object storage across
// multiple providers including local filesystem, MinIO, and Volcengine TOS.
//
// # Supported Providers
//
//   - local: Local filesystem storage
//   - minio: MinIO and S3-compatible storage
//   - volcengine: Volcengine TOS (Toutiao Object Storage)
//
// # Basic Usage
//
// Create a storage client:
//
//	import "github.com/poly-workshop/go-webmods/object_storage"
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderLocal,
//	    ProviderConfig: object_storage.ProviderConfig{
//	        BasePath: "/data/uploads",
//	    },
//	})
//	if err != nil {
//	    panic(err)
//	}
//
// # Local Filesystem Provider
//
// Store files on the local filesystem:
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderLocal,
//	    ProviderConfig: object_storage.ProviderConfig{
//	        BasePath: "/var/data",  // Base directory for all files
//	    },
//	})
//
// # MinIO / S3-Compatible Provider
//
// Connect to MinIO or any S3-compatible service:
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderMinio,
//	    ProviderConfig: object_storage.ProviderConfig{
//	        Endpoint:  "minio.example.com:9000",
//	        Region:    "us-east-1",
//	        AccessKey: "minioadmin",
//	        SecretKey: "minioadmin",
//	        Bucket:    "mybucket",
//	        BasePath:  "uploads/",  // Optional prefix for all keys
//	    },
//	})
//
// # Volcengine TOS Provider
//
// Connect to Volcengine TOS:
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderVolcengine,
//	    ProviderConfig: object_storage.ProviderConfig{
//	        Endpoint:    "tos-cn-beijing.volces.com",
//	        Region:      "cn-beijing",
//	        AccessKey:   "your-access-key",
//	        SecretKey:   "your-secret-key",
//	        Bucket:      "mybucket",
//	        UseInternal: false,  // Use internal endpoint for in-region access
//	    },
//	})
//
// # Common Operations
//
// Save a file:
//
//	file, _ := os.Open("photo.jpg")
//	defer file.Close()
//	size, err := storage.Save("photos/photo.jpg", file)
//
// List files in a directory:
//
//	objects, err := storage.List("photos/")
//	for _, obj := range objects {
//	    info, _ := obj.Stat()
//	    fmt.Printf("%s - %d bytes\n", info.Name(), info.Size())
//	}
//
// Open and read a file:
//
//	obj, err := storage.Open("photos/photo.jpg")
//	if err != nil {
//	    panic(err)
//	}
//	defer obj.Close()
//
//	data, err := io.ReadAll(obj)
//	// Or use obj.Read(), obj.Seek() for streaming
//
// Get file metadata:
//
//	info, err := storage.Stat("photos/photo.jpg")
//	fmt.Printf("Name: %s, Size: %d, Modified: %s\n",
//	    info.Name(), info.Size(), info.ModTime())
//
// Delete a file:
//
//	err := storage.Delete("photos/photo.jpg")
//
// # Interface Design
//
// The ObjectStorage interface provides a consistent API across all providers:
//
//	type ObjectStorage interface {
//	    Save(path string, r io.Reader) (int64, error)
//	    List(path string) ([]Object, error)
//	    Open(path string) (Object, error)
//	    Stat(path string) (os.FileInfo, error)
//	    Delete(path string) error
//	}
//
// Objects returned by Open() and List() implement io.ReadSeekCloser:
//
//	type Object interface {
//	    io.ReadSeekCloser
//	    Stat() (os.FileInfo, error)
//	}
//
// This allows seamless use with Go's standard I/O functions.
//
// # Configuration with Viper
//
// Example configuration file (configs/default.yaml):
//
//	object_storage:
//	  provider: local
//	  base_path: /var/data/uploads
//
// For production (configs/production.yaml):
//
//	object_storage:
//	  provider: minio
//	  endpoint: minio.prod.example.com:9000
//	  region: us-east-1
//	  access_key: ${MINIO_ACCESS_KEY}
//	  secret_key: ${MINIO_SECRET_KEY}
//	  bucket: prod-uploads
//
// Loading configuration:
//
//	import (
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/object_storage"
//	)
//
//	app.Init(".")
//	cfg := app.Config()
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderType(cfg.GetString("object_storage.provider")),
//	    ProviderConfig: object_storage.ProviderConfig{
//	        Endpoint:  cfg.GetString("object_storage.endpoint"),
//	        Region:    cfg.GetString("object_storage.region"),
//	        AccessKey: cfg.GetString("object_storage.access_key"),
//	        SecretKey: cfg.GetString("object_storage.secret_key"),
//	        Bucket:    cfg.GetString("object_storage.bucket"),
//	        BasePath:  cfg.GetString("object_storage.base_path"),
//	    },
//	})
//
// # Best Practices
//
//   - Use local provider for development and testing
//   - Use MinIO or cloud providers for production
//   - Always close Object instances after use (defer obj.Close())
//   - Use BasePath to organize files and support multi-tenancy
//   - Validate file sizes before calling Save() to prevent abuse
//   - Use Stat() to check file existence before Open()
//   - Handle errors appropriately (file not found, permission denied, etc.)
//   - For large files, use streaming (Open/Read) instead of loading into memory
//   - Consider using signed URLs for direct client uploads (provider-specific)
//
// # Error Handling
//
// Common errors:
//   - File not found: Open() and Stat() return errors
//   - Permission denied: All operations may return permission errors
//   - Invalid configuration: NewObjectStorage() returns errors
//   - Network errors: MinIO and TOS operations may fail on network issues
//
// Always check returned errors:
//
//	obj, err := storage.Open("file.txt")
//	if err != nil {
//	    if errors.Is(err, os.ErrNotExist) {
//	        // Handle file not found
//	    }
//	    return err
//	}
//	defer obj.Close()
//
// # Provider-Specific Notes
//
// Local Provider:
//   - Automatically creates directories as needed
//   - BasePath is the root directory
//   - Paths are relative to BasePath
//
// MinIO Provider:
//   - Supports any S3-compatible service
//   - BasePath is used as a key prefix
//   - Requires bucket to exist before use
//
// Volcengine TOS Provider:
//   - UseInternal=true uses internal endpoint (for in-region VMs)
//   - Supports Volcengine-specific features
//   - Region must match bucket region
package object_storage
