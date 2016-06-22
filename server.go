package sdk

import (
	"laatoo/framework/core/server"
)

func Server(config string) error {
	return server.Main(config)
}
