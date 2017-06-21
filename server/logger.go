package server

import (
	"laatoo/sdk/core"
	"laatoo/sdk/log"
)

type Logger interface {
	core.ServerElement
	log.Logger
}
