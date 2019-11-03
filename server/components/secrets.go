package components

import (
	"laatoo/sdk/server/core"
)

type SecretsManager interface {
	Get(ctx core.ServerContext, key string) ([]byte, bool)
	Put(ctx core.ServerContext, key string, val []byte) error
}
