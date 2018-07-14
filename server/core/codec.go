package core

import (
	"laatoo/sdk/server/ctx"
)

type Codec interface {
	Unmarshal(ctx.Context, []byte, interface{}) error
	Marshal(ctx.Context, interface{}) ([]byte, error)
}
