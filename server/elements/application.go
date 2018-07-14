package elements

import (
	"laatoo/sdk/server/core"
)

type Application interface {
	core.ServerElement
	GetApplet(name string) (Applet, bool)
}
