package elements

import (
	"laatoo/sdk/server/core"
)

type Applet interface {
	core.ServerElement
	ServerElementHandle
}
