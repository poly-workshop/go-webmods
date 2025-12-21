package object_storage

import (
	"io"
	"os"
)

type Object interface {
	io.ReadSeekCloser
	Stat() (os.FileInfo, error)
}

type ObjectStorage interface {
	Save(path string, r io.Reader) (int64, error)
	List(path string) ([]Object, error)
	Open(path string) (Object, error)
	Stat(path string) (os.FileInfo, error)
	Delete(path string) error
}
