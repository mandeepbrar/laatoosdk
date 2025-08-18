package datatypes

import "laatoo.io/sdk/ctx"

type Serializable interface {
	ReadAll(ctx.Context, Codec, SerializableReader) error
	WriteAll(ctx.Context, Codec, SerializableWriter) error
}
