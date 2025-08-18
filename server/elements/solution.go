package elements

import (
	"laatoo.io/sdk/server/core"
)

type Solution interface {
	core.ServerElement
	GetPeers(filter string) ([]string, error)
}
