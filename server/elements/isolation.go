package elements

import (
	"laatoo.io/sdk/server/core"
)

type Isolation interface {
	core.ServerElement
	GetIsolationId() string
	//GetApplet(name string) (Applet, bool)
}
