package data

import (
	"laatoo/sdk/core"
	"laatoo/sdk/log"
	"time"
)

///auditable entities must have UpdatedBy and UpdatedOn fields to support auditing through update queries
type Auditable interface {
	Storable
	IsNew() bool
	SetUpdatedOn(string)
	SetUpdatedBy(string)
	SetCreatedBy(string)
	GetCreatedBy() string
}

func Audit(ctx core.RequestContext, item interface{}) {
	if item != nil {
		auditable, aok := item.(Auditable)
		if aok {
			usr := ctx.GetUser()
			if usr != nil {
				id := usr.GetId()
				if auditable.IsNew() {
					auditable.SetCreatedBy(id)
				}
				auditable.SetUpdatedBy(id)
				auditable.SetUpdatedOn(time.Now().Format(time.UnixDate))
			} else {
				log.Logger.Info(ctx, "Could not audit entity. User nil")
			}
		} else {
			updateMap, mapok := item.(map[string]interface{})
			if mapok {
				usr := ctx.GetUser()
				if usr != nil {
					id := usr.GetId()
					updateMap["UpdatedBy"] = id
					updateMap["UpdatedOn"] = time.Now().Format(time.UnixDate)
				} else {
					log.Logger.Info(ctx, "Could not audit map. User nil")
				}
			}
		}
	}
}
