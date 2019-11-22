package core

import (
	"io"
	"laatoo/sdk/server/ctx"
	"time"
)

type SerializableReader interface {
	io.Reader
	Start() error
	Bytes() []byte
	ReadProp(ctx ctx.Context, cdc Codec, prop string) (SerializableReader, bool, error)
	ReadBytes(ctx ctx.Context, cdc Codec, prop string) ([]byte, bool, error)
	ReadInt(ctx ctx.Context, cdc Codec, prop string, val *int) (bool, error)
	ReadInt64(ctx ctx.Context, cdc Codec, prop string, val *int64) (bool, error)
	ReadString(ctx ctx.Context, cdc Codec, prop string, val *string) (bool, error)
	ReadFloat32(ctx ctx.Context, cdc Codec, prop string, val *float32) (bool, error)
	ReadFloat64(ctx ctx.Context, cdc Codec, prop string, val *float64) (bool, error)
	ReadBool(ctx ctx.Context, cdc Codec, prop string, val *bool) (bool, error)
	ReadObject(ctx ctx.Context, cdc Codec, prop string, val interface{}) (bool, error)
	ReadMap(ctx ctx.Context, cdc Codec, prop string, val *map[string]interface{}) (bool, error)
	ReadArray(ctx ctx.Context, cdc Codec, prop string, val interface{}) (bool, error)
	ReadTime(ctx ctx.Context, cdc Codec, prop string, val *time.Time) (bool, error)
}
