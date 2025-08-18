package elements

import (
	"laatoo.io/sdk/server/core"
)

type Application interface {
	core.ServerElement
	//GetApplet(name string) (Applet, bool)
}
