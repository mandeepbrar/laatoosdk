package components

import (
	"io"
	"time"

	"laatoo.io/sdk/server/core"
)

// StorageComponent is the platform's file/object storage contract.
//
// THREE BACKENDS IMPLEMENT IT AND THEY DIVERGE ON MOST METHODS —
// laatoomodules/storage/dev/plugins/{filesystemstorage,googlestorage,s3storage}. A `bucket` is a
// directory on one and a real object-store bucket on the other two, an empty bucket name resolves
// to the backend's configured default on the read/write methods (but NOT on CreateBucket or
// DeleteBucket), and several methods mean genuinely different things per backend. Read the method
// doc before assuming portability; the ones most likely to surprise are ListFiles, GetFullPath,
// DeleteFiles and DeleteBucket.
type StorageComponent interface {
	// Open returns a reader for the named object. The caller must Close it.
	//
	// The error is the backend's own not-found — os.ErrNotExist, storage.ErrObjectNotExist, a minio
	// error — and is not normalised, so never match on its text.
	//
	// googlestorage builds a FRESH GCS CLIENT PER CALL and never closes it: the defer is commented
	// out (googlestorage.go:158-166), so a hot read path leaks connections on that backend.
	Open(ctx core.RequestContext, bucket, fileName string) (io.ReadCloser, error)

	// OpenForWrite returns a writer for the named object, creating or truncating it. The caller
	// must Close it — on the object stores nothing is durable until Close returns without error.
	//
	// NOT IMPLEMENTED ON s3storage: it returns errors.NotImplemented unconditionally
	// (s3storage.go:134-141). Use SaveFile if the deployment may be backed by S3 or minio.
	//
	// On googlestorage an object written this way does NOT receive the public-read ACL that
	// CreateFile applies (googlestorage.go:175 versus :144), so the same object is world-readable
	// or not depending purely on which method wrote it.
	OpenForWrite(ctx core.RequestContext, bucket, fileName string) (io.WriteCloser, error)

	// SaveFile writes the whole of inpStr to the named object and returns the stored file name.
	//
	// The returned string is fileName echoed back — NOT a URL and not a backend-assigned id. A
	// caller that needs a location must call GetFullPath or SignedReadURL afterwards.
	//
	// It CLOSES inpStr on the success path on filesystemstorage and googlestorage, and does not
	// close it on the error paths; s3storage hands the reader to minio and never closes it at all.
	// Defer your own Close rather than relying on this.
	//
	// On googlestorage this takes the CreateFile path and therefore inherits its public-read ACL —
	// see CreateFile.
	SaveFile(ctx core.RequestContext, bucket string, inpStr io.ReadCloser, fileName string, contentType string) (string, error)

	// GetFullPath returns a location string for the object. It is pure string computation: nothing
	// is contacted and the object need not exist.
	//
	// WHAT IT RETURNS IS NOT THE SAME KIND OF THING ON EACH BACKEND, so it is not portable and must
	// never be handed to a browser as a link without knowing which backend is configured:
	//   - filesystemstorage: an absolute path on THIS SERVER's filesystem (filesystem.go:123-125)
	//   - googlestorage: an https URL — storage.googleapis.com when the module is `public`,
	//     otherwise storage.cloud.google.com, which needs a signed-in Google identity
	//     (googlestorage.go:199-204)
	//   - s3storage: an https URL when `public`, and an EMPTY STRING otherwise, with no error and
	//     no way to tell that apart from a legitimately empty answer (s3storage.go:160-167)
	//
	// For a link a client can follow against a private object, use SignedReadURL.
	GetFullPath(ctx core.RequestContext, bucket string, fileName string) string

	// ServeFile writes the object into the CURRENT REQUEST's response rather than returning it: it
	// sets the response on ctx and returns nil on success.
	//
	// The response it sets differs by backend. filesystemstorage streams the bytes with a
	// Content-Type guessed from the FILENAME EXTENSION rather than from anything stored
	// (filesystem.go:85-88); s3storage streams with no content type at all; googlestorage streams
	// too, EXCEPT when its module is configured `public: true`, where it instead sets a REDIRECT to
	// the public GCS URL and never reads the object (googlestorage.go:179-182). A caller cannot
	// assume the bytes flow through this server.
	ServeFile(ctx core.RequestContext, bucket string, fileName string) error

	// CreateFile returns a writer for a new object, recording contentType where the backend
	// supports it. The caller must Close it.
	//
	// contentType is honoured by googlestorage and s3storage and IGNORED by filesystemstorage,
	// which has nowhere to record it (filesystem.go:52-57).
	//
	// googlestorage STAMPS EVERY OBJECT CREATED THIS WAY WITH A PUBLIC-READ ACL
	// (AllUsers/RoleReader, googlestorage.go:144) regardless of the module's `public` setting.
	// Anything written through CreateFile or SaveFile on that backend is readable by anyone holding
	// the URL, whether or not the deployment intended public files.
	//
	// On s3storage this DEADLOCKS as written: it builds an io.Pipe and calls SaveFile on the read
	// half synchronously before returning the write half (s3storage.go:105-116), so the upload
	// blocks waiting for bytes only the not-yet-returned writer could supply.
	CreateFile(ctx core.RequestContext, bucket string, fileName string, contentType string) (io.WriteCloser, error)

	// CopyFile streams the object into dest.
	//
	// dest is NOT closed — the caller owns it (laatoomodules/storage/dev/plugins/common/common.go:
	// 51-63) — and is left partially written if the copy fails midway, so a caller writing into a
	// destination object must decide for itself whether to keep what landed.
	CopyFile(ctx core.RequestContext, bucket string, fileName string, dest io.WriteCloser) error

	// ListFiles lists objects in the bucket, filtered by pattern.
	//
	// PATTERN MEANS THREE DIFFERENT THINGS, AND ONE BACKEND IGNORES IT:
	//   - filesystemstorage: a filepath.Glob pattern, returning FULL FILESYSTEM PATHS rather than
	//     object names (filesystem.go:144-147)
	//   - s3storage: a key PREFIX handed to ListObjectsV2, non-recursive, returning object keys
	//     (s3storage.go:177-194)
	//   - googlestorage: IGNORED ENTIRELY — it lists the bucket with a nil query and returns every
	//     object name it sees (googlestorage.go:232-256)
	//
	// googlestorage additionally TRUNCATES SILENTLY: it sizes the result from
	// it.PageInfo().Remaining(), which is what is buffered in the current page and not the object
	// count, so a bucket larger than one page comes back short with no error.
	//
	// s3storage passes bucket through WITHOUT the default-bucket resolution every other method on
	// that backend applies (s3storage.go:186), so an empty bucket name lists nothing there.
	ListFiles(ctx core.RequestContext, bucket string, pattern string) ([]string, error)

	// DeleteFiles removes the SINGLE named object. Despite the plural name it takes one file name
	// and accepts no pattern.
	//
	// THE BOOL AND THE ERROR SPLIT BY BACKEND for the same call. Deleting an absent file on
	// filesystemstorage is (false, nil) — an idempotent no-op (filesystem.go:98-108). On
	// googlestorage and s3storage an absent object is the backend's not-found ERROR, returned as
	// (false, err). A delete-if-present loop that is correct on one backend fails on the others.
	DeleteFiles(ctx core.RequestContext, bucket string, fileName string) (bool, error)

	// Exists reports whether the object is present.
	//
	// A NEGATIVE ANSWER IS NOT PROOF OF ABSENCE. All three backends collapse every error into
	// false, so wrong credentials, an unreachable bucket and a genuine miss are the same answer
	// (filesystem.go:59-66, googlestorage.go:148-157, s3storage.go:118-124). Where the difference
	// matters, call Open and read the error.
	//
	// googlestorage defers client.Close() BEFORE checking the client-creation error
	// (googlestorage.go:149-150), so on that backend this panics rather than returning false when
	// the client cannot be built.
	Exists(ctx core.RequestContext, bucket string, fileName string) bool

	// CreateBucket creates the bucket.
	//
	// NOT IDEMPOTENT on any backend: creating one that already exists is an error (os.Mkdir's
	// EEXIST, GCS's conflict, minio's BucketAlreadyOwnedByYou). filesystemstorage uses os.Mkdir
	// rather than MkdirAll, so it also fails when the storage root does not yet exist
	// (filesystem.go:149-152).
	//
	// bucket is used RAW on filesystemstorage and s3storage, without the default-bucket
	// substitution the read and write methods apply — an empty name there is an empty name, not the
	// default bucket. googlestorage DOES substitute it (googlestorage.go:206-217), so the same
	// empty argument creates the default bucket on one backend and fails on the other two.
	CreateBucket(ctx core.RequestContext, bucket string) error

	// DeleteBucket removes the bucket.
	//
	// HOW DESTRUCTIVE IT IS DIFFERS. filesystemstorage does os.RemoveAll: it deletes the directory
	// AND EVERYTHING UNDER IT, and returns nil when the directory does not exist
	// (filesystem.go:154-157). googlestorage and s3storage delegate to the object store, which
	// refuses to delete a bucket that still holds objects. The same call is a recursive wipe on one
	// backend and a refusal on the others.
	//
	// bucket is used raw on filesystemstorage and s3storage; see CreateBucket.
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
