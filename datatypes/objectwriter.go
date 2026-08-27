package datatypes

import (
	"io"
	"time"

	"laatoo.io/sdk/ctx"
)

// SerializableWriter is a cursor that builds one encoded object, handed to a type's generated
// WriteAll so it can emit its own fields by name. A codec creates it; a plugin normally meets it
// only inside generated code or when embedding a nested object by hand.
//
// # The call order is part of the format, not a convention
//
// The writer emits the encoding as it goes, so the delimiters come from the calls themselves:
//
//	Start()          // emits the opening delimiter
//	Write…(…)        // emits each property
//	Close()          // emits the CLOSING delimiter
//	Bytes()          // now valid
//
// Bytes does NOT close the object. Reading it before Close yields output truncated just before
// the final delimiter — on JSON, a document with no closing brace. That parses nowhere and the
// writer reports no error, so the failure appears wherever the bytes are later read.
//
// Writing before Start is the mirror image: the properties land outside the object.
//
// # Every Write dereferences its argument
//
// The scalar methods take a pointer and dereference it unconditionally. A nil pointer PANICS —
// it is not treated as "omit this property". Skip the call for a field that should not be
// written.
type SerializableWriter interface {
	io.WriteCloser

	// Start emits the opening of the object. Call it before any Write.
	Start() error

	// Bytes returns what has been written so far. Call Close first: this does not close the
	// object, and the bytes are incomplete until it is closed.
	Bytes() []byte

	// WriteBytes writes a binary property. Panics on a nil pointer.
	WriteBytes(ctx ctx.Context, cdc Codec, prop string, val *[]byte) error

	// WriteInt writes an integer property. Panics on a nil pointer.
	WriteInt(ctx ctx.Context, cdc Codec, prop string, val *int) error

	// WriteInt32 writes a 32-bit integer property. Panics on a nil pointer.
	WriteInt32(ctx ctx.Context, cdc Codec, prop string, val *int32) error

	// WriteInt64 writes a 64-bit integer property. Panics on a nil pointer.
	WriteInt64(ctx ctx.Context, cdc Codec, prop string, val *int64) error

	// WriteString writes a string property. Panics on a nil pointer.
	WriteString(ctx ctx.Context, cdc Codec, prop string, val *string) error

	// WriteFloat32 writes a 32-bit float property. Panics on a nil pointer.
	WriteFloat32(ctx ctx.Context, cdc Codec, prop string, val *float32) error

	// WriteFloat64 writes a 64-bit float property. Panics on a nil pointer.
	WriteFloat64(ctx ctx.Context, cdc Codec, prop string, val *float64) error

	// WriteBool writes a boolean property. Panics on a nil pointer.
	WriteBool(ctx ctx.Context, cdc Codec, prop string, val *bool) error

	// WriteObject writes a nested object property, encoding val with the codec.
	WriteObject(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// WriteMap writes a mapping property.
	WriteMap(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// WriteArray writes a list property.
	WriteArray(ctx ctx.Context, cdc Codec, prop string, val interface{}) error

	// WriteTime writes a timestamp property. Panics on a nil pointer.
	WriteTime(ctx ctx.Context, cdc Codec, prop string, val *time.Time) error
}
