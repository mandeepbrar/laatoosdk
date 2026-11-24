package components

import (
	"laatoo.io/sdk/server/core"
)

// SecretsManager is a backing store for secret values, implemented by a plugin service
// (vaultsecrets, gcpsecrets) or by the server's built-in encrypted disk store.
//
// Which one a level uses is configuration: the server resolves the service named by the
// `secretsvc` setting and asserts it to this interface, and falls back to the built-in store when
// nothing is named (laatooserver/src/core/secretsmanager.go:23-47). Managers are then chained
// solution -> application -> isolation, each level consulting its parent on a miss
// (secretsmanager.go:49-85) — which is why the (value, ok, error) triple below is a contract and
// not a convention. It is codified, with the reasoning, in
// laatoomodules/secrets/dev/plugins/common/common.go:17-46 (Hit/Miss/Failure).
type SecretsManager interface {
	// Get resolves key, returning (value, true, nil) on a hit, (nil, false, nil) on a MISS, and
	// (nil, false, err) on a FAILURE — the store was unreachable, the credential was rejected, or
	// the stored data was not the shape this backend writes.
	//
	// The distinction is load-bearing and the two failure modes run in opposite directions.
	//
	// Reporting an absent key as an error ABORTS THE WHOLE LOOKUP CHAIN: the server returns on the
	// first non-nil error and never consults the parent manager (secretsmanager.go:52-55), so a
	// solution-level secret stops being inherited by an application that has its own store. The
	// server's own built-in disk store does exactly this — it wraps any storer error, including
	// "not found", into an error (laatooserver/src/core/secretsmanagerimpl.go:68-89).
	//
	// Reporting a genuine failure as a miss is worse: the lookup falls through to the parent
	// store, and the caller either gets a value from the wrong level or a NotFound naming the
	// manager rather than the real cause, with nothing pointing at the denied permission or the
	// corrupt secret.
	//
	// Two further behaviours of the built-in store that this interface does not express: it
	// appends "_test" to the key when the server runs in test mode (secretsmanagerimpl.go, Get),
	// so the same key resolves to a different secret there; and it performs no key validation,
	// unlike the plugin backends, which reject empty keys, path separators and traversal segments
	// before the key becomes a path segment
	// (laatoomodules/secrets/dev/plugins/common/common.go:48-66).
	Get(ctx core.ServerContext, key string) ([]byte, bool, error)

	// Put stores val under key.
	//
	// How the server's own layers write (laatooserver/src/core):
	//
	//   - the built-in disk store encrypts val with the seed-password key and writes it in the
	//     format `laatoo security addsecret` writes (secretsmanagerimpl.go, Put);
	//   - the directory of secrets delivered to a deployed server is read-only, and a Put that
	//     reaches it is refused with an error (secretsfilestore.go);
	//   - the aggregating manager writes to the level's own store, or, at a level with no store of
	//     its own, to the nearest ancestor that has one, which is the store Get reads there
	//     (secretsmanager.go, put). A write is dropped only when no level up to the root has a
	//     store.
	//
	// Plugin backends write as their stores do: vaultsecrets creates a new KV v2 version, and
	// gcpsecrets adds a version, creating the secret first only when its `createmissing` setting
	// is on. A caller that must know the value is readable reads the key back.
	//
	// Implementations must return an error only for a real failure, and must not treat overwriting
	// an existing key as one unless the backend genuinely forbids it.
	Put(ctx core.ServerContext, key string, val []byte) error
}
