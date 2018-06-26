package components

import (
	"io"
	"laatoo/sdk/server/core"
)

type StorageComponent interface {
	Open(ctx core.RequestContext, fileName string) (io.ReadCloser, error)
	SaveFile(ctx core.RequestContext, inpStr io.ReadCloser, fileName string, contentType string) (string, error)
	GetFullPath(ctx core.RequestContext, fileName string) string
	ServeFile(ctx core.RequestContext, fileName string) error
	CreateFile(ctx core.RequestContext, fileName string, contentType string) (io.WriteCloser, error)
	ListFiles(ctx core.RequestContext, pattern string) ([]string, error)
	Exists(ctx core.RequestContext, fileName string) bool
}
