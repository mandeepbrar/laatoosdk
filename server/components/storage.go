package components

import (
	"io"
	"time"

	"laatoo.io/sdk/server/core"
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

	// SignedReadURL returns a URL that grants read access to one object for ttl, and
	// SignedWriteURL returns one that grants write access. They let a caller hand a
	// client a credential for a single object instead of proxying the bytes itself.
	//
	// The returned string is a complete URL with the signature already inside it —
	// there is no separate token to attach. A caller must pass it on unmodified:
	// appending a query parameter, changing the port, or normalising the path
	// invalidates the signature.
	//
	// ttl belongs to the caller, not the backend. Whoever issues the credential is the
	// only party that knows how long the operation it authorizes should stay possible,
	// so no default is applied here.
	//
	// Authorization is not delegated by issuing one of these. The caller decides whether
	// the request should be allowed and issues a URL only afterwards; the storage
	// backend verifies its own signature and never sees who is on the other end.
	//
	// A backend with no signing scheme — a plain filesystem has none — returns
	// errors.NotImplemented and an empty string, so callers can tell "this backend
	// cannot issue credentials" from "issuing failed". Returning an unsigned URL
	// instead would look like success and fail only at the point of use.
	SignedReadURL(ctx core.RequestContext, bucket string, fileName string, ttl time.Duration) (string, error)
	SignedWriteURL(ctx core.RequestContext, bucket string, fileName string, ttl time.Duration) (string, error)
}
