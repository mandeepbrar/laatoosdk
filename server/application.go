package server

import (
	"laatoo/sdk/core"
)

type Application interface {
	core.ServerElement
	GetApplet(name string) (Applet, bool)
}
