package server

import (
	"laatoo/sdk/core"
)

type Applet interface {
	core.ServerElement
	ServerElementHandle
}
