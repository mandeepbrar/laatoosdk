package components

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// CacheComponent is a key/value cache partitioned into named buckets, obtained from
// elements.CacheManager.GetCache. Each implementation documents how it meets this contract,
// including what it shares across replicas; write against the one the deployment configures.
//
// No method reports a cache miss as an error. The (interface{}, bool) methods report it in the
// bool; the map-returning methods report it by what the map holds (see GetMulti).
type CacheComponent interface {
	// PutTempObject stores item under key in bucket, expiring ttl from now. A ttl <= 0 stores it
	// with no expiry, as PutObject does. An implementation may remove an expired value only when it
	// is next read, so until then it can still appear in ListKeys.
	PutTempObject(ctx core.RequestContext, bucket string, key string, item interface{}, ttl time.Duration) error

	// PutObject stores item under key in bucket, with no expiry.
	//
	// item is serialised with the component's configured codec, so it must round-trip through that
	// codec: GetObject and GetIntoObject unmarshal into a freshly created object of the type the
	// caller names, never handing back the value passed in here. Do not mutate item after storing
	// it; an implementation without a codec may keep it by reference.
	PutObject(ctx core.RequestContext, bucket string, key string, item interface{}) error

	// PutObjects stores every key/value pair in vals, exactly as PutObject would. It is not atomic:
	// on an error some pairs may already be stored, and which ones is unspecified.
	PutObjects(ctx core.RequestContext, bucket string, vals utils.StringMap) error

	// GetObject reads key, creates a new instance of objectType through the server context,
	// unmarshals the stored bytes into it and returns it. objectType is the registered object name
	// ("plugin.EntityName") -- the same string ctx.CreateObject takes.
	//
	// false means no usable value came back: a miss, an unregistered objectType and a failed
	// unmarshal all read the same, with no error. Reading back false immediately after a
	// successful write is a serialisation or object-registration problem, not an eviction.
	GetObject(ctx core.RequestContext, bucket string, key string, objectType string) (interface{}, bool)

	// GetIntoObject reads key and unmarshals it into obj, which the caller allocates.
	//
	// A miss is an error here, and its kind is not part of the contract: never branch on it. Use
	// Get or GetObject when absent has to be told apart from broken.
	GetIntoObject(ctx core.RequestContext, bucket string, key string, obj interface{}) error

	// Get returns the raw stored value for key, and whether it was present.
	//
	// The raw value is the implementation's stored form -- normally the codec-encoded []byte, or a
	// []byte item exactly as it was written -- and never the object PutObject was handed.
	// Type-assert defensively, or use GetObject and GetIntoObject, which decode. It is the accessor
	// for values written as []byte and for the counters Increment and Decrement maintain.
	Get(ctx core.RequestContext, bucket string, key string) (interface{}, bool)

	// GetObjects reads several keys and returns a map of key to a decoded objectType instance. A
	// key that is absent or fails to decode is either left out of the map or present with a nil
	// value: test the value, not the presence of the key. There is no error return, so a failure
	// of the whole read looks like every key being absent.
	GetObjects(ctx core.RequestContext, bucket string, keys []string, objectType string) utils.StringMap

	// GetMulti reads several keys and returns a map of key to raw stored value, in the form Get
	// returns it. As with GetObjects, a missing key is either left out or present with a nil
	// value, and a failure of the whole read looks like every key being absent.
	GetMulti(ctx core.RequestContext, bucket string, keys []string) utils.StringMap

	// Delete removes key from bucket. Deleting a key that is not there is not an error, so a nil
	// return does not mean anything was deleted.
	Delete(ctx core.RequestContext, bucket string, key string) error

	// Increment adds one to the integer stored at key, creating it at 1 when it does not exist.
	//
	// It does not return the new value: read it back with Get, in the implementation's own
	// counter encoding. A counter is not readable through GetObject or GetIntoObject, and
	// incrementing a key PutObject wrote is undefined. Atomic across every caller that shares the
	// cache.
	Increment(ctx core.RequestContext, bucket string, key string) error

	// Decrement subtracts one from the integer stored at key, creating it at -1 when it does not
	// exist. The read-back and atomicity terms are Increment's.
	Decrement(ctx core.RequestContext, bucket string, key string) error

	// ListKeys returns the keys held in bucket, in no particular order. It may still list a key
	// whose value has expired but not yet been removed, so it can disagree with Get.
	ListKeys(ctx core.RequestContext, bucket string) ([]string, error)

	// CreateTempObject stores item at key in bucket, expiring ttl from now, ONLY if the key holds
	// no live value, and answers whether this call stored it. It is what a single-use token needs:
	//
	//	first, err := cache.CreateTempObject(ctx, "redeemedcodes", codeId, []byte("1"), time.Until(codeExpiry))
	//	if err != nil { ... } // undecided -- never treat it as claimed
	//	if !first { ... }     // someone presented this code before
	//
	// No combination of the other methods gives it: PutTempObject expires but is not atomic -- two
	// callers can both read "absent" and both write -- and Increment is atomic but never expires.
	//
	// Of any number of concurrent callers sharing the cache, exactly one gets true. ttl must be
	// positive; a ttl <= 0 is refused with an error. A key whose value has expired may be claimed
	// again. An error means the claim could not be decided (the store was unreachable, or
	// contention outlasted the retries).
	CreateTempObject(ctx core.RequestContext, bucket string, key string, item interface{}, ttl time.Duration) (bool, error)
}
