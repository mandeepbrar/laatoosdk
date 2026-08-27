package datatypes

import "laatoo.io/sdk/ctx"

// Serializable is the platform's own field-by-field encoding contract, used instead of reflection
// over struct tags. Entity implementations are GENERATED: the codegen emits one Read*/Write* call
// per declared field, then chains the same call into each embedded info struct (StorageInfo,
// TrackingInfo, DeletionInfo, TenantInfo).
//
// Because the methods are an explicit list of named fields, a field the implementation does not
// name is simply not transported — no error is raised on either side. Hand-writing an
// implementation and forgetting a field therefore loses that field silently.
type Serializable interface {
	// ReadAll populates the receiver from rdr, one named field at a time, and returns on the
	// first reader error.
	//
	// It is ADDITIVE, not a reset: a field absent from the input is left holding whatever the
	// receiver already had, so decoding into a reused object merges rather than replaces. It also
	// stops at the first failure, which leaves the receiver partially populated — a caller that
	// keeps the object after an error is holding a half-decoded record.
	ReadAll(ctx.Context, Codec, SerializableReader) error

	// WriteAll emits the receiver's fields into wtr, returning on the first writer error.
	//
	// CAUTION: two of the codec entry points discard this error. BinaryCodec.MarshalWriter
	// (laatooserver/src/codecs/binarycodec.go:93) and JsonCodec.MarshalWriter
	// (laatooserver/src/codecs/jsoncodec.go:149) both call WriteAll without checking it and then
	// return (wtr.Bytes(), nil), so a failed encode reaches the caller as a SUCCESS carrying
	// truncated bytes. The Marshal/Unmarshal helpers in codecs/misc.go do check it; MarshalWriter
	// is the path that does not.
	WriteAll(ctx.Context, Codec, SerializableWriter) error
}
