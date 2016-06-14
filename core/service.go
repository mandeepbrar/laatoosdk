package core

import (
	"laatoo/sdk/config"
)

type ServiceFunc func(ctx RequestContext) error

type Service interface {
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Invoke(RequestContext) error
}
