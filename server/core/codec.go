package core

import (
	"io"
	"laatoo/sdk/server/ctx"
)

type Codec interface {
	Unmarshal(ctx.Context, []byte, interface{}) error
	Marshal(ctx.Context, interface{}) ([]byte, error)
	Encode(c ctx.Context, outStream io.Writer, val interface{}) error
	Decode(c ctx.Context, inpStream io.Reader, val interface{}) error
	UnmarshalSerializable(ctx.Context, []byte, Serializable) error
	MarshalSerializable(ctx.Context, Serializable) ([]byte, error)
	UnmarshalSerializableProps(ctx.Context, []byte, Serializable, map[string]interface{}) error
	MarshalSerializableProps(ctx.Context, Serializable, map[string]interface{}) ([]byte, error)
	UnmarshalReader(ctx.Context, SerializableReader, Serializable) error
	MarshalWriter(ctx.Context, SerializableWriter, Serializable) ([]byte, error)
	UnmarshalSerializableFromStream(c ctx.Context, rdr io.Reader, obj Serializable) error
	UnmarshalSerializableFromStreamProps(c ctx.Context, rdr io.Reader, obj Serializable, props map[string]interface{}) error
	MarshalSerializableToStream(c ctx.Context, wtr io.Writer, obj Serializable) error
	MarshalSerializableToStreamProps(c ctx.Context, wtr io.Writer, obj Serializable, props map[string]interface{}) error
}
