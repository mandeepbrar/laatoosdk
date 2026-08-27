package components

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// CacheComponent is a key/value cache partitioned into named buckets, obtained from
// elements.CacheManager.GetCache.
//
// THREE BACKENDS IMPLEMENT THIS AND THEY DO NOT AGREE. Write against the one the deployment
// configures, and check the method's doc before relying on a behaviour:
//
//   - the server's built-in NATS JetStream KeyValue cache (laatooserver/src/core/natscache.go) —
//     shared across replicas, and the only one that is still a cache on more than one pod.
//   - the memorycache plugin (laatoomodules/cache/dev/plugins/memorycache) — a process-local map.
//     On more than one replica each pod holds its OWN copy: a value written on one is a miss on
//     the next, and Delete on one does not evict the others. It is a per-process memo, not a
//     cache, and nothing at runtime says so.
//   - the rediscache plugin (laatoomodules/cache/dev/plugins/rediscache) — shared, but several of
//     its methods do not work as written; see PutTempObject and ListKeys.
//
// A bucket is a namespace, and how real it is depends on the backend: NATS gives each bucket its
// own JetStream KV store, memorycache its own map, redis only a key prefix.
//
// No method reports a cache miss as an error. The (interface{}, bool) methods report it in the
// bool; the map-returning methods report it by what the map holds — and that differs by backend,
// see GetMulti.
type CacheComponent interface {
	// PutTempObject stores item under key in bucket with a per-key expiry of ttl.
	//
	// ONLY THE NATS BACKEND HONOURS ttl. memorycache accepts the argument and DROPS IT ENTIRELY,
	// delegating straight to PutObject (memorycacheservice.go:59-62); its only expiry is a
	// bucket-wide `ttl` fixed in module config at startup. The value is stored and readable — it
	// simply never expires on the schedule the caller asked for, and no error says so.
	//
	// NATS enforces expiry LAZILY ON READ, tagging the value with an absolute deadline rather than
	// using a bucket TTL: an expired key is removed by the first Get or GetIntoObject that sees it.
	// Until something reads it, it still occupies the store and is still returned by ListKeys.
	// A ttl <= 0 is treated as PutObject, with no expiry at all.
	//
	// BROKEN ON REDIS: it issues SETEX with the value and the seconds transposed
	// (rediscacheservice.go:57 — redis expects SETEX key seconds value), which redis rejects, so
	// every call returns an error and nothing is stored.
	PutTempObject(ctx core.RequestContext, bucket string, key string, item interface{}, ttl time.Duration) error

	// PutObject stores item under key in bucket, with no expiry.
	//
	// item is serialised with the component's configured codec, so it must round-trip through that
	// codec: GetObject and GetIntoObject unmarshal into a freshly created object of the type the
	// caller names, never handing back the value passed in here. NATS stores a []byte item verbatim
	// without encoding it, which is also what makes Get's return type backend-specific.
	//
	// On memorycache with no `encoding` configured the value is stored BY REFERENCE and not copied
	// (memorycacheservice.go:66-68), so a later mutation of item is visible to every reader.
	//
	// On redis a module-level `ttl` setting silently turns this into an expiring write
	// (rediscacheservice.go:72-76) — which then hits the same transposed-SETEX defect as
	// PutTempObject, so configuring that setting makes every PutObject fail.
	PutObject(ctx core.RequestContext, bucket string, key string, item interface{}) error

	// PutObjects stores every key/value pair in vals, exactly as PutObject would.
	//
	// NOT ATOMIC on any backend, and not all-or-nothing: NATS and memorycache write one key at a
	// time and return on the first failure, leaving everything already written in place
	// (natscache.go:210-218). Go map iteration order is random, so which subset survives a partial
	// failure is not predictable and not reproducible.
	PutObjects(ctx core.RequestContext, bucket string, vals utils.StringMap) error

	// GetObject reads key, creates a new instance of objectType through the server context,
	// unmarshals the stored bytes into it and returns it. The bool reports whether a usable value
	// came back.
	//
	// FALSE ALSO MEANS "FOUND BUT UNUSABLE". Every backend collapses a miss, an unregistered
	// objectType and a failed unmarshal into (nil, false) with no error anywhere
	// (natscache.go:265-282, rediscacheservice.go:119-143, memorycacheservice.go:120-138). Reading
	// back false immediately after a successful write is a serialisation or object-registration
	// problem, not an eviction, and nothing will say so.
	//
	// objectType is the registered object name ("plugin.EntityName") — the same string
	// ctx.CreateObject takes.
	//
	// On memorycache with no codec configured this degrades to Get and IGNORES objectType.
	GetObject(ctx core.RequestContext, bucket string, key string, objectType string) (interface{}, bool)

	// GetIntoObject reads key and unmarshals it into obj, which the caller allocates.
	//
	// THE ERROR TELLS YOU NOTHING PORTABLE ABOUT A MISS. NATS wraps the store's not-found
	// (natscache.go:226-229), memorycache returns errors.NotFound (memorycacheservice.go:104), and
	// redis returns an InternalError reading "data value is not byte array" — because a redis miss
	// is a nil reply rather than an error (rediscacheservice.go:210-213). Never branch on the error
	// kind here; use Get or GetObject when absent has to be told apart from broken.
	//
	// On NATS an expired-but-not-yet-collected key returns a CACHE_KEY_EXPIRED error and is deleted.
	//
	// On memorycache with no `encoding` configured this ALWAYS fails with an internal error
	// (memorycacheservice.go:96-98) — the same configuration on which PutObject and Get work fine.
	GetIntoObject(ctx core.RequestContext, bucket string, key string, obj interface{}) error

	// Get returns the raw stored value for key, and whether it was present.
	//
	// WHAT COMES BACK DIFFERS BY BACKEND, and on none of them is it the object PutObject was handed.
	// NATS returns the stored []byte with any expiry header stripped (natscache.go:244-262); redis
	// returns redigo's reply, normally []byte; memorycache returns the encoded []byte, or — with no
	// codec configured — the original value by reference. Type-assert defensively, or use GetObject
	// and GetIntoObject, which do the decoding.
	//
	// This is the accessor for values written as []byte, and the only one for the counters
	// Increment and Decrement maintain — NATS writes those as eight raw big-endian bytes
	// (natscache.go:331-333) and redis in its own integer encoding, so reading a counter back is
	// necessarily backend-specific.
	Get(ctx core.RequestContext, bucket string, key string) (interface{}, bool)

	// GetObjects reads several keys and returns a map of key to a decoded objectType instance.
	//
	// THE MAP'S TREATMENT OF A MISS IS BACKEND-SPECIFIC, the same split GetMulti has: memorycache
	// and redis insert an explicit nil value for every key that was absent or failed to decode,
	// while NATS omits the key altogether (natscache.go:297-306, memorycacheservice.go:154-177,
	// rediscacheservice.go:172-202). `v, ok := m[k]` therefore yields ok==true with a nil v on two
	// of the three. Test the value, not the presence of the key.
	//
	// redis returns a NIL MAP (not an empty one) when the MGET itself fails or an object cannot be
	// created; the others return whatever they managed to gather. No backend returns an error, so a
	// total backend failure is indistinguishable from every key being absent.
	GetObjects(ctx core.RequestContext, bucket string, keys []string, objectType string) utils.StringMap

	// GetMulti reads several keys and returns a map of key to raw stored value, in the form Get
	// would return it.
	//
	// THE MISS CONVENTION SPLITS THE BACKENDS. NATS omits absent keys from the map entirely
	// (natscache.go:285-294); memorycache and redis insert the key with a nil value
	// (memorycacheservice.go:140-152, rediscacheservice.go:157-168). Code that decides presence
	// with the two-value map lookup is correct on one backend and silently wrong on the other two —
	// check for a nil value as well.
	//
	// There is no error return: a backend-level failure yields an empty map, which is exactly what
	// every key being absent looks like.
	GetMulti(ctx core.RequestContext, bucket string, keys []string) utils.StringMap

	// Delete removes key from bucket.
	//
	// Deleting a key that is not there is NOT an error on any backend — NATS explicitly swallows
	// not-found (natscache.go:118-120) and redis DEL reports zero removals as success. A nil return
	// therefore does not mean anything was deleted.
	Delete(ctx core.RequestContext, bucket string, key string) error

	// Increment adds one to the integer stored at key, creating it at 1 when it does not exist.
	//
	// IT DOES NOT RETURN THE NEW VALUE. Read it back with Get, in the backend's own encoding: NATS
	// writes eight raw big-endian bytes and accepts either that or a decimal string on read
	// (natscache.go:309-343), redis uses INCR. A counter is NOT readable through GetObject or
	// GetIntoObject, which expect codec-encoded bytes — and incrementing a key that PutObject wrote
	// reinterprets that value's bytes as a number rather than failing.
	//
	// Atomic on NATS (compare-and-swap, five attempts, then a CAS_FAILED error) and on redis. On
	// memorycache it is atomic only within the one process, which is the caveat on that whole
	// backend.
	Increment(ctx core.RequestContext, bucket string, key string) error

	// Decrement subtracts one from the integer stored at key. See Increment — the encoding,
	// atomicity and read-back caveats are identical, and decrementing an absent key creates it at
	// -1 rather than failing.
	Decrement(ctx core.RequestContext, bucket string, key string) error

	// ListKeys returns the keys currently held in bucket.
	//
	// EXPECT IT TO DISAGREE WITH Get. On NATS expiry is enforced lazily on read and this applies no
	// such check (natscache.go:383-400), so an expired key stays listed until something reads it.
	//
	// ON REDIS THE RETURNED STRINGS ARE NOT USABLE KEYS. Redis keys are built with
	// fmt.Sprintf("%s_%#v", bucket, variants) (laatoomodules/cache/dev/plugins/common/utils.go:5-7),
	// so a key stored as "abc" lives at `bucket_[]interface {}{"abc"}`; ListKeys strips only the
	// `bucket_` prefix and hands back the Go-syntax remainder (rediscacheservice.go:295-330).
	// Feeding one back into Get finds nothing.
	//
	// Ordering is unspecified everywhere, and memorycache returns nil rather than an empty slice
	// for an empty bucket.
	ListKeys(ctx core.RequestContext, bucket string) ([]string, error)
}
