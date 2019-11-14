package components

import (
	"io"
	"laatoo/sdk/server/core"
)

type StorageComponent interface {
	Open(ctx core.RequestContext, bucket, fileName string) (io.ReadCloser, error)
	OpenForWrite(ctx core.RequestContext, bucket, fileName string) (io.WriteCloser, error)
	SaveFile(ctx core.RequestContext, bucket string, inpStr io.ReadCloser, fileName string, contentType string) (string, error)
	GetFullPath(ctx core.RequestContext, bucket string, fileName string) string
	ServeFile(ctx core.RequestContext, bucket string, fileName string) error
	CreateFile(ctx core.RequestContext, bucket string, fileName string, contentType string) (io.WriteCloser, error)
	CopyFile(ctx core.RequestContext, bucket string, fileName string, dest io.WriteCloser) error
	ListFiles(ctx core.RequestContext, bucket string, pattern string) ([]string, error)
	DeleteFiles(ctx core.RequestContext, bucket string, fileName string) (bool, error)
	Exists(ctx core.RequestContext, bucket string, fileName string) bool
	CreateBucket(ctx core.RequestContext, bucket string) error
	DeleteBucket(ctx core.RequestContext, bucket string) error
}
