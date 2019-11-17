package core

import (
	"laatoo/sdk/server/ctx"
	"time"
)

type SerializableReader interface {
	Read(ctx ctx.Context, cdc Codec, prop string) (SerializableReader, error)
	ReadBytes(ctx ctx.Context, cdc Codec) ([]byte, error)
	ReadByte(ctx ctx.Context, cdc Codec, prop string, val *byte) error
	ReadInt(ctx ctx.Context, cdc Codec, prop string, val *int) error
	ReadInt64(ctx ctx.Context, cdc Codec, prop string, val *int64) error
	ReadString(ctx ctx.Context, cdc Codec, prop string, val *string) error
	ReadFloat32(ctx ctx.Context, cdc Codec, prop string, val *float32) error
	ReadFloat64(ctx ctx.Context, cdc Codec, prop string, val *float64) error
	ReadBool(ctx ctx.Context, cdc Codec, prop string, val *bool) error
	ReadObject(ctx ctx.Context, cdc Codec, prop, objtype string) (SerializableReader, error)
	ReadMap(ctx ctx.Context, cdc Codec, prop string, val *map[string]interface{}) error
	ReadArray(ctx ctx.Context, cdc Codec, prop string, objtype string, val interface{}) error
	ReadTime(ctx ctx.Context, cdc Codec, prop string, val *time.Time) error
}
