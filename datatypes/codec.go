package datatypes

import (
	"io"

	"laatoo.io/sdk/ctx"
)

// Codec converts values to and from an encoding. Obtain one from a context with
// GetCodec(encoding) rather than constructing it.
//
// # Which encodings exist
//
// The platform resolves "json", "multipart" and "bin". Anything else — including "protobuf",
// whose codec type exists but is not wired into the lookup — returns not-found rather than an
// error, so check the second return before using the result.
//
// # One instance per encoding, shared by every caller
//
// A codec is created on first use of its encoding and cached for the life of the process, so
// every request handler holds the SAME instance. An implementation must therefore be safe for
// concurrent use and must not carry per-call state on itself.
//
// # Three families of method, and which to reach for
//
//   - Marshal/Unmarshal and Encode/Decode work on any value by reflection, and are what most
//     code wants.
//   - The Serializable family drives a type's own generated ReadAll/WriteAll, which is faster and
//     is what the platform uses internally for entities.
//   - The Reader/Writer variants let a caller supply the cursor, for embedding one object's
//     encoding inside another's.
//
// # The Props variants do not do what their name suggests
//
// Every method taking a props map ACCEPTS AND DISCARDS IT: the property map is never consulted,
// so each Props variant behaves exactly like its plain counterpart. They cannot be used to
// select, filter or project fields. Treat them as aliases until that changes.
type Codec interface {
	// Unmarshal decodes bytes into val by reflection.
	//
	// Empty input, and input that is exactly the JSON literal null, are treated as "nothing to
	// do": the call returns nil and leaves val untouched rather than clearing it.
	Unmarshal(ctx.Context, []byte, interface{}) error

	// Marshal encodes any value to bytes by reflection. A nil value returns (nil, nil) rather
	// than an encoded null.
	Marshal(ctx.Context, interface{}) ([]byte, error)

	// Encode writes any value to a stream, for a payload that should not be assembled in memory
	// first. A nil value writes nothing and returns nil.
	Encode(c ctx.Context, outStream io.Writer, val interface{}) error

	// Decode reads a value from a stream. A nil target returns nil without consuming the stream.
	Decode(c ctx.Context, inpStream io.Reader, val interface{}) error

	// UnmarshalSerializable decodes into a type that implements Serializable, driving its own
	// ReadAll rather than reflecting over its fields.
	//
	// Absent properties are not errors — see SerializableReader — so a payload missing every
	// field decodes successfully and leaves the object as it was.
	UnmarshalSerializable(ctx.Context, []byte, Serializable) error

	// MarshalSerializable encodes a Serializable by driving its WriteAll.
	MarshalSerializable(ctx.Context, Serializable) ([]byte, error)

	// UnmarshalSerializableProps decodes a Serializable. The props map is accepted and ignored;
	// this is UnmarshalSerializable.
	UnmarshalSerializableProps(ctx.Context, []byte, Serializable, map[string]interface{}) error

	// MarshalSerializableProps encodes a Serializable. The props map is accepted and ignored.
	MarshalSerializableProps(ctx.Context, Serializable, map[string]interface{}) ([]byte, error)

	// UnmarshalReader decodes an object from a reader the caller already holds, for reading a
	// nested object out of an enclosing document.
	UnmarshalReader(ctx.Context, SerializableReader, Serializable) error

	// MarshalWriter encodes an object into a writer the caller already holds, for embedding one
	// object inside another's encoding.
	//
	// It reports no failure: the error from the object's WriteAll is discarded and the returned
	// error is always nil, so a partially written object is indistinguishable from a complete
	// one. Inspect the bytes if it matters.
	MarshalWriter(ctx.Context, SerializableWriter, Serializable) ([]byte, error)

	// UnmarshalSerializableFromStream decodes a Serializable directly from a stream, without
	// buffering the whole payload.
	UnmarshalSerializableFromStream(c ctx.Context, rdr io.Reader, obj Serializable) error

	// UnmarshalSerializableFromStreamProps decodes from a stream. The props map is ignored.
	UnmarshalSerializableFromStreamProps(c ctx.Context, rdr io.Reader, obj Serializable, props map[string]interface{}) error

	// MarshalSerializableToStream encodes a Serializable directly to a stream.
	MarshalSerializableToStream(c ctx.Context, wtr io.Writer, obj Serializable) error

	// MarshalSerializableToStreamProps encodes to a stream. The props map is ignored.
	MarshalSerializableToStreamProps(c ctx.Context, wtr io.Writer, obj Serializable, props map[string]interface{}) error
}
