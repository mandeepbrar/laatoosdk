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
	// appends "_test" to the key when the server runs in test mode (secretsmanagerimpl.go:69-72),
	// so the same key resolves to a different secret there; and it performs no key validation,
	// unlike the plugin backends, which reject empty keys, path separators and traversal segments
	// before the key becomes a path segment
	// (laatoomodules/secrets/dev/plugins/common/common.go:48-66).
	Get(ctx core.ServerContext, key string) ([]byte, bool, error)

	// Put stores val under key.
	//
	// A NIL RETURN DOES NOT MEAN THE SECRET WAS STORED. Two layers discard writes silently:
	//
	//   - the server's built-in disk store implements Put as `return nil` and writes nothing
	//     (laatooserver/src/core/secretsmanagerimpl.go:91-93) — secrets in that store are placed
	//     out of band, by `laatoo security addsecret`;
	//   - the server's aggregating manager returns nil when the level has no store of its own
	//     (secretsmanager.go:87-95), and unlike Get it never delegates to the parent, so a write
	//     at a level with no store is dropped rather than inherited upward.
	//
	// A Put that must actually persist therefore requires a backend that genuinely writes —
	// vaultsecrets (VaultSecretsSvc.Put, which creates a new KV v2 version) or gcpsecrets — and
	// the caller cannot tell the difference from the return value alone. Read the key back if the
	// write has to be confirmed.
	//
	// Implementations must return an error only for a real failure, and must not treat overwriting
	// an existing key as one unless the backend genuinely forbids it.
	Put(ctx core.ServerContext, key string, val []byte) error
}
