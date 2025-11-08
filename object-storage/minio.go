package object_storage

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioObjectStorage struct {
	ctx    context.Context
	client *minio.Client

	bucket   string
	basePath string
}

type MinioObject struct {
	ctx        context.Context
	client     *minio.Client
	bucket     string
	key        string
	rangeStart int64
}

func (o *MinioObject) Close() error {
	o.client = nil
	return nil
}

func (o *MinioObject) Seek(offset int64, whence int) (int64, error) {
	stat, err := o.client.StatObject(o.ctx, o.bucket, o.key, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}
	objectSize := stat.Size

	switch whence {
	case io.SeekStart:
		o.rangeStart = offset
	case io.SeekCurrent:
		o.rangeStart += offset
	case io.SeekEnd:
		o.rangeStart = objectSize - 1 - offset
	}
	if o.rangeStart < 0 {
		o.rangeStart = 0
	}
	if o.rangeStart >= objectSize {
		o.rangeStart = objectSize - 1
	}
	return o.rangeStart, nil
}

func (o *MinioObject) Stat() (os.FileInfo, error) {
	stat, err := o.client.StatObject(o.ctx, o.bucket, o.key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return FileInfo{
		path:         o.key,
		name:         path.Base(o.key),
		objectSize:   stat.Size,
		lastModified: stat.LastModified,
	}, nil
}

func (o *MinioObject) Read(p []byte) (n int, err error) {
	l := int64(len(p))

	// Set range options for partial read
	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(o.rangeStart, o.rangeStart+l-1); err != nil {
		return 0, err
	}

	obj, err := o.client.GetObject(o.ctx, o.bucket, o.key, opts)
	if err != nil {
		stat, statErr := o.Stat()
		if statErr != nil {
			return 0, statErr
		}
		if o.rangeStart >= stat.Size() {
			return 0, io.EOF
		}
		return 0, err
	}
	defer func() { _ = obj.Close() }()

	n, err = obj.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	}
	o.rangeStart += int64(n)
	return n, err
}

func NewMinioObjectStorage(config ProviderConfig) (*MinioObjectStorage, error) {
	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: true, // Use HTTPS by default
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}

	return &MinioObjectStorage{
		ctx:      context.Background(),
		client:   client,
		bucket:   config.Bucket,
		basePath: config.BasePath,
	}, nil
}

func (s *MinioObjectStorage) getObjectPath(objectPath string) string {
	return path.Join(s.basePath, objectPath)
}

func (s *MinioObjectStorage) Save(path string, r io.Reader) (int64, error) {
	objectPath := s.getObjectPath(path)

	// Upload object to MinIO
	info, err := s.client.PutObject(s.ctx, s.bucket, objectPath, r, -1, minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}

func (s *MinioObjectStorage) List(path string) ([]Object, error) {
	prefix := s.getObjectPath(path)

	// List objects with prefix
	objects := []Object{}
	for objInfo := range s.client.ListObjects(s.ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if objInfo.Err != nil {
			return nil, objInfo.Err
		}

		// Skip directories (objects ending with /)
		if objInfo.Key[len(objInfo.Key)-1] == '/' {
			continue
		}

		objects = append(objects, &MinioObject{
			ctx:    s.ctx,
			client: s.client,
			bucket: s.bucket,
			key:    objInfo.Key,
		})
	}

	return objects, nil
}

func (s *MinioObjectStorage) Open(path string) (Object, error) {
	objectPath := s.getObjectPath(path)

	return &MinioObject{
		ctx:    s.ctx,
		client: s.client,
		bucket: s.bucket,
		key:    objectPath,
	}, nil
}

func (s *MinioObjectStorage) Stat(path string) (os.FileInfo, error) {
	objectPath := s.getObjectPath(path)

	stat, err := s.client.StatObject(s.ctx, s.bucket, objectPath, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		path:         objectPath,
		name:         path,
		objectSize:   stat.Size,
		lastModified: stat.LastModified,
	}, nil
}

func (s *MinioObjectStorage) Delete(path string) error {
	objectPath := s.getObjectPath(path)

	return s.client.RemoveObject(s.ctx, s.bucket, objectPath, minio.RemoveObjectOptions{})
}
