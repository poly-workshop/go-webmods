package object_storage

import (
	"context"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

type TOSObjectStorage struct {
	ctx            context.Context
	client         *tos.ClientV2
	internalClient *tos.ClientV2
	useInternal    bool

	bucket   string
	basePath string
}

type TOSObject struct {
	ctx        context.Context
	client     *tos.ClientV2
	bucket     string
	key        string
	rangeStart int64
}

func (o *TOSObject) Close() error {
	o.client = nil
	return nil
}

func (o *TOSObject) Seek(offset int64, whence int) (int64, error) {
	resp, err := o.client.HeadObjectV2(o.ctx, &tos.HeadObjectV2Input{
		Bucket: o.bucket,
		Key:    o.key,
	})
	if err != nil {
		return 0, err
	}
	objectSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

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

func (o *TOSObject) Stat() (os.FileInfo, error) {
	resp, err := o.client.HeadObjectV2(o.ctx, &tos.HeadObjectV2Input{
		Bucket: o.bucket,
		Key:    o.key,
	})
	if err != nil {
		return nil, err
	}
	objectSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return FileInfo{
		path:         o.key,
		name:         path.Base(o.key),
		objectSize:   objectSize,
		lastModified: lastModified,
	}, nil
}

func (o *TOSObject) Read(p []byte) (n int, err error) {
	l := int64(len(p))
	resp, err := o.client.GetObjectV2(o.ctx, &tos.GetObjectV2Input{
		Bucket:     o.bucket,
		Key:        o.key,
		RangeStart: o.rangeStart,
		RangeEnd:   o.rangeStart + l - 1,
	})
	if err != nil {
		state, err := o.Stat()
		if err != nil {
			return 0, err
		}
		if o.rangeStart >= state.Size() {
			return 0, io.EOF
		}
		return 0, err
	}
	defer func() { _ = resp.Content.Close() }()

	body, err := io.ReadAll(resp.Content)
	if err != nil {
		return 0, err
	}
	n = copy(p, body)
	o.rangeStart += int64(n)
	return n, nil
}

func NewTOSObjectStorage(config ProviderConfig) (*TOSObjectStorage, error) {
	client, err := tos.NewClientV2(
		config.Endpoint,
		tos.WithRegion(config.Region),
		tos.WithCredentials(tos.NewStaticCredentials(config.AccessKey, config.SecretKey)),
	)
	if err != nil {
		return nil, err
	}

	config.Endpoint = strings.ReplaceAll(config.Endpoint, "volces.com", "ivolces.com")
	internalClient, err := tos.NewClientV2(
		config.Endpoint,
		tos.WithRegion(config.Region),
		tos.WithCredentials(tos.NewStaticCredentials(config.AccessKey, config.SecretKey)),
	)
	if err != nil {
		return nil, err
	}

	return &TOSObjectStorage{
		ctx:            context.Background(),
		client:         client,
		internalClient: internalClient,
		useInternal:    config.UseInternal,

		bucket:   config.Bucket,
		basePath: config.BasePath,
	}, nil
}

func (s *TOSObjectStorage) getObjectPath(objectPath string) string {
	return path.Join(s.basePath, objectPath)
}

func (s *TOSObjectStorage) Save(path string, r io.Reader) (int64, error) {
	client := s.client
	if s.useInternal {
		client = s.internalClient
	}

	_, err := client.PutObjectV2(s.ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: s.bucket,
			Key:    s.getObjectPath(path),
		},
		Content: r,
	})
	if err != nil {
		return 0, err
	}
	fileInfo, err := s.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func (s *TOSObjectStorage) List(path string) ([]Object, error) {
	client := s.client
	if s.useInternal {
		client = s.internalClient
	}
	res, err := client.ListObjectsType2(s.ctx, &tos.ListObjectsType2Input{
		Bucket: s.bucket,
		Prefix: s.getObjectPath(path),
	})
	if err != nil {
		return nil, err
	}
	objs := []Object{}
	for _, file := range res.Contents {
		objs = append(objs, &TOSObject{
			ctx:    s.ctx,
			client: client,
			bucket: s.bucket,
			key:    file.Key,
		})
	}
	return objs, nil
}

func (s *TOSObjectStorage) Open(path string) (Object, error) {
	client := s.client
	if s.useInternal {
		client = s.internalClient
	}

	return &TOSObject{
		ctx:    s.ctx,
		client: client,
		bucket: s.bucket,
		key:    s.getObjectPath(path),
	}, nil
}

func (s *TOSObjectStorage) Stat(path string) (os.FileInfo, error) {
	client := s.client
	if s.useInternal {
		client = s.internalClient
	}

	resp, err := client.HeadObjectV2(s.ctx, &tos.HeadObjectV2Input{
		Bucket: s.bucket,
		Key:    s.getObjectPath(path),
	})
	if err != nil {
		return nil, err
	}
	objectSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	fileInfo := &FileInfo{
		path:         s.getObjectPath(path),
		name:         path,
		objectSize:   objectSize,
		lastModified: lastModified,
	}
	return fileInfo, nil
}

func (s *TOSObjectStorage) Delete(path string) error {
	client := s.client
	if s.useInternal {
		client = s.internalClient
	}

	_, err := client.DeleteObjectV2(s.ctx, &tos.DeleteObjectV2Input{
		Bucket: s.bucket,
		Key:    s.getObjectPath(path),
	})
	return err
}
