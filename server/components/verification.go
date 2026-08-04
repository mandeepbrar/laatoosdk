package components

// KeySource names where the material for verifying a credential comes from.
type KeySource string

const (
	// KeySourceNamed resolves a key from the secret store by name. The platform, or whoever
	// registered the caller, holds it.
	KeySourceNamed KeySource = "named"
	// KeySourceExternal fetches a key set published by an external issuer. The platform holds no
	// key at all and the credential names which key of the set signed it.
	KeySourceExternal KeySource = "external"
)

// VerificationDescriptor is how a service declares the way an inbound machine credential presented
// on its channel is verified. The service declares; the security manager verifies. No service
// implements credential cryptography, and no service holds key material.
//
// Two shapes are expressible, and the second is the reason this is a descriptor rather than a key:
// a credential minted by an external issuer is verified against a key set that rotates on that
// issuer's schedule and is selected per credential, which a single key cannot express.
type VerificationDescriptor struct {
	// Source selects between a key this platform holds and a key set an issuer publishes.
	Source KeySource

	// KeyName is the secret-store entry holding the verification key. Set when Source is
	// KeySourceNamed. It is a name, never key bytes — resolution happens in the security manager
	// through components.SecretsManager so key material never reaches service state.
	KeyName string

	// KeySetLocation is where the issuer publishes its key set. Set when Source is
	// KeySourceExternal.
	KeySetLocation string

	// Issuer is the expected issuer of the credential. A credential asserting a different issuer is
	// refused.
	Issuer string

	// Audience is the expected audience of the credential. A credential asserting a different
	// audience is refused. Without this, a credential minted for one endpoint is valid at every
	// endpoint that trusts the same issuer.
	Audience string

	// Algorithms are the signing algorithms accepted for this credential, pinned here rather than
	// read from the credential itself.
	//
	// This is a security property, not a convenience. A verifier that takes the algorithm from the
	// credential's own header accepts a credential signed with a published verification key used as
	// a shared secret, and accepts one declaring no algorithm at all. Pinning out of band closes
	// both by construction. An empty list means the descriptor is incomplete and verification must
	// refuse rather than fall back to a default.
	Algorithms []string
}
