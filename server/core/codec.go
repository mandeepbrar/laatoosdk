package core

import (
	"laatoo/sdk/server/ctx"
)

type Codec interface {
	Unmarshal(ctx.Context, []byte, interface{}) error
	Marshal(ctx.Context, interface{}) ([]byte, error)
	UnmarshalSerializable(ctx.Context, []byte, Serializable) error
	MarshalSerializable(ctx.Context, Serializable) ([]byte, error)
	UnmarshalSerializableProps(ctx.Context, []byte, Serializable, map[string]interface{}) error
	MarshalSerializableProps(ctx.Context, Serializable, map[string]interface{}) ([]byte, error)
	UnmarshalReader(ctx.Context, SerializableReader, Serializable) error
	MarshalWriter(ctx.Context, SerializableWriter) ([]byte, error)
}
