package data

import "laatoo/sdk/core"

//Object stored by data service
type CustomLoader interface {
	LoadData(ctx core.RequestContext, services ...DataComponent) error
}
