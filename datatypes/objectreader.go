package datatypes

import (
	"io"
	"time"

	"laatoo.io/sdk/ctx"
)

// SerializableReader is a cursor over one encoded object, handed to a type's generated ReadAll so
// it can pull its own fields out by name. A codec creates it; a plugin normally meets it only
// inside generated code or when decoding a nested object by hand.
//
// # An absent property is not an error, and that is the thing to know
//
// Every Read method returns nil when the property is missing and leaves the target UNTOUCHED. It
// does not zero it, and it does not report the absence. So:
//
//   - A payload missing every field decodes successfully and the object keeps whatever it held.
//     Decoding into a reused object therefore merges rather than replaces.
//   - A field cannot be told apart from one explicitly present at its zero value.
//
// ReadProp and ReadBytes are the exceptions that can signal absence: they return a nil result
// with a nil error, so test the result rather than the error.
//
// Errors are reserved for a property that IS present but cannot be converted to the requested
// type.
type SerializableReader interface {
	io.Reader

	// Start prepares the cursor. Call it before reading; on some encodings it is a no-op.
	Start() error

	// Bytes returns the encoding of the node the cursor is on, for handing a nested value to
	// something that decodes bytes.
	Bytes() []byte

	// ReadProp returns a cursor over a nested object, for descending into it. A nil reader with
	// a nil error means the property is absent.
	ReadProp(ctx ctx.Context, cdc Codec, prop string) (SerializableReader, error)

	// ReadBytes reads a binary property. A nil result with a nil error means absent.
	ReadBytes(ctx ctx.Context, cdc Codec, prop string) ([]byte, error)

	// ReadInt reads an integer property into val, leaving it untouched when absent.
	ReadInt(ctx ctx.Context, cdc Codec, prop string, val *int) error

	// ReadInt32 reads a 32-bit integer property into val, leaving it untouched when absent.
	ReadInt32(ctx ctx.Context, cdc Codec, prop string, val *int32) error

	// ReadInt64 reads a 64-bit integer property into val, leaving it untouched when absent.
	ReadInt64(ctx ctx.Context, cdc Codec, prop string, val *int64) error

	// ReadString reads a string property into val, leaving it untouched when absent.
	ReadString(ctx ctx.Context, cdc Codec, prop string, val *string) error

	// ReadFloat32 reads a 32-bit float property into val, leaving it untouched when absent.
	ReadFloat32(ctx ctx.Context, cdc Codec, prop string, val *float32) error

	// ReadFloat64 reads a 64-bit float property into val, leaving it untouched when absent.
	ReadFloat64(ctx ctx.Context, cdc Codec, prop string, val *float64) error

	// ReadBool reads a boolean property into val, leaving it untouched when absent.
	ReadBool(ctx ctx.Context, cdc Codec, prop string, val *bool) error

	// ReadObject decodes a nested object into val, which must be a pointer to the target.
	ReadObject(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// ReadMap decodes a mapping property into val.
	ReadMap(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// ReadArray decodes a list property into val.
	ReadArray(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// ReadTime reads a timestamp property into val, leaving it untouched when absent.
	ReadTime(ctx ctx.Context, cdc Codec, prop string, val *time.Time) error
}
