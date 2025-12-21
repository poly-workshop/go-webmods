package object_storage

import (
	"io"
	"os"
	"path/filepath"
)

// LocalObjectStorage implements ObjectStorage interface for local file system
type LocalObjectStorage struct {
	basePath string
}

// LocalObject implements Object interface for local files
type LocalObject struct {
	file *os.File
}

// NewLocalObjectStorage creates a new local object storage instance
func NewLocalObjectStorage(config ProviderConfig) (*LocalObjectStorage, error) {
	basePath := config.BasePath
	if basePath == "" {
		basePath = "./data"
	}

	// Ensure the base directory exists
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, err
	}

	return &LocalObjectStorage{
		basePath: basePath,
	}, nil
}

// getObjectPath returns the full path for an object
func (s *LocalObjectStorage) getObjectPath(objectPath string) string {
	return filepath.Join(s.basePath, objectPath)
}

// Save saves a file to the local storage
func (s *LocalObjectStorage) Save(path string, r io.Reader) (int64, error) {
	fullPath := s.getObjectPath(path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return 0, err
	}

	// Create and write to file
	file, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close() // Ignore close error as we've already written successfully
	}()

	written, err := io.Copy(file, r)
	if err != nil {
		return 0, err
	}

	return written, nil
}

// List lists objects in the given path
func (s *LocalObjectStorage) List(path string) ([]Object, error) {
	fullPath := s.getObjectPath(path)

	// Check if path exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return []Object{}, nil
	}

	var objects []Object

	err := filepath.Walk(fullPath, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and the root path itself
		if info.IsDir() || walkPath == fullPath {
			return nil
		}

		file, err := os.Open(walkPath)
		if err != nil {
			return err
		}

		objects = append(objects, &LocalObject{file: file})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// Open opens a file for reading
func (s *LocalObjectStorage) Open(path string) (Object, error) {
	fullPath := s.getObjectPath(path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return &LocalObject{file: file}, nil
}

// Stat returns file information
func (s *LocalObjectStorage) Stat(path string) (os.FileInfo, error) {
	fullPath := s.getObjectPath(path)
	return os.Stat(fullPath)
}

// Delete deletes a file
func (s *LocalObjectStorage) Delete(path string) error {
	fullPath := s.getObjectPath(path)
	return os.Remove(fullPath)
}

// LocalObject methods

// Read implements io.Reader
func (o *LocalObject) Read(p []byte) (n int, err error) {
	return o.file.Read(p)
}

// Seek implements io.Seeker
func (o *LocalObject) Seek(offset int64, whence int) (int64, error) {
	return o.file.Seek(offset, whence)
}

// Close implements io.Closer
func (o *LocalObject) Close() error {
	return o.file.Close()
}

// Stat returns file information
func (o *LocalObject) Stat() (os.FileInfo, error) {
	return o.file.Stat()
}
