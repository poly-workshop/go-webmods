package object_storage

import (
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	path         string
	name         string
	objectSize   int64
	lastModified time.Time
}

func (i FileInfo) Name() string {
	return i.name
}

func (i FileInfo) Size() int64 {
	return i.objectSize
}

func (i FileInfo) ModTime() time.Time {
	return i.lastModified
}

func (i FileInfo) IsDir() bool {
	return strings.HasSuffix(i.path, "/")
}

func (i FileInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (i FileInfo) Sys() any {
	return nil
}
