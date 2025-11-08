package object_storage

import "fmt"

type ProviderType string

const (
	ProviderLocal      ProviderType = "local"
	ProviderVolcengine ProviderType = "volcengine"
	ProviderMinio      ProviderType = "minio"
)

type Config struct {
	ProviderType
	ProviderConfig
}

type ProviderConfig struct {
	Endpoint    string
	Region      string
	AccessKey   string
	SecretKey   string
	Bucket      string
	BasePath    string
	UseInternal bool
}

func NewObjectStorage(cfg Config) (ObjectStorage, error) {
	switch cfg.ProviderType {
	case ProviderLocal:
		return NewLocalObjectStorage(cfg.ProviderConfig)
	case ProviderVolcengine:
		return NewTOSObjectStorage(cfg.ProviderConfig)
	case ProviderMinio:
		return NewMinioObjectStorage(cfg.ProviderConfig)
	default:
		return nil, fmt.Errorf("unsupported object storage provider: %s", cfg.ProviderType)
	}
}
