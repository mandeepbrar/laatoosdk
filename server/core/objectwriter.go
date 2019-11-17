package core

import (
	"laatoo/sdk/server/ctx"
	"time"
)

type SerializableWriter interface {
	Bytes() []byte
	WriteBytes(ctx ctx.Context, cdc Codec, prop string, val *[]byte) error
	WriteInt(ctx ctx.Context, cdc Codec, prop string, val *int) error
	WriteInt64(ctx ctx.Context, cdc Codec, prop string, val *int64) error
	WriteString(ctx ctx.Context, cdc Codec, prop string, val *string) error
	WriteFloat32(ctx ctx.Context, cdc Codec, prop string, val *float32) error
	WriteFloat64(ctx ctx.Context, cdc Codec, prop string, val *float64) error
	WriteBool(ctx ctx.Context, cdc Codec, prop string, val *bool) error
	WriteObject(ctx ctx.Context, cdc Codec, prop string, val interface{}) error
	WriteMap(ctx ctx.Context, cdc Codec, prop string, val *map[string]interface{}) error
	WriteArray(ctx ctx.Context, cdc Codec, prop string, val interface{}) error
	WriteTime(ctx ctx.Context, cdc Codec, prop string, val *time.Time) error
}
