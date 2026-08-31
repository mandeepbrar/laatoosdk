package data

import "laatoo.io/sdk/server/core"

const (
	CONF_DATA_SVCS                = "dataservices"
	CONF_BASE_SVC                 = "baseservice"
	CONF_DATA_OBJECT              = "object"
	CONF_DATA_POSTSAVE            = "postsave"
	CONF_DATA_WORKFLOW_ENABLED    = "workflow"
	CONF_DATA_POSTLOAD            = "postload"
	CONF_DATA_MULTITENANT         = "multitenant"
	CONF_DATA_POSTUPDATE          = "postupdate"
	CONF_DATA_PRESAVE             = "presave"
	CONF_DATA_CACHEABLE           = "cacheable"
	CONF_DATA_AUDITABLE           = "auditable"
	CONF_DATA_REFOPS              = "refops"
	CONF_DATA_COLLECTION          = "collection"
	CONF_DATA_EMBEDDED_DOC_SEARCH = "embedded_doc_search"
	CONF_DATA_SOFTDELETE          = "softdelete"
	CONF_PRESAVE_MSG              = "storable_presave"
	CONF_POSTSAVE_MSG             = "storable_postsave"
	CONF_PREUPDATE_MSG            = "storable_preupdate"
	CONF_POSTUPDATE_MSG           = "storable_postupdate"
	CONF_NEWOBJ_MSG               = "storable_new"
	// Delete notifications. They complete the lifecycle: presave/postsave and
	// preupdate/postupdate already exist, and delete was the one mutation a rule could not
	// observe -- so anything deriving state from a record (an edge row, a projection, a cache
	// entry) had no way to learn the record was gone, and silently kept the derived state.
	//
	// These travel the SYNCHRONOUS rules path, ctx.SendSynchronousMessage, exactly as the save
	// and update pairs do. That is deliberate and is the whole point: the call returns an error,
	// so a listener that cannot remove its derived state FAILS the delete rather than leaving
	// the record gone and the derivation behind. The asynchronous data-event path
	// (EmitDataEvent) cannot express that -- it is fire-and-forget and returns nothing.
	CONF_PREDELETE_MSG  = "storable_predelete"
	CONF_POSTDELETE_MSG = "storable_postdelete"

	// CONF_DEFAULT_DATACONNECTION is the SOLUTION-level key naming which dataconnection a lookup
	// resolves to when it names none. The DataManager reads it in Initialize and REFUSES TO BOOT
	// without it: a connection is a factory instance name, so there is no value the platform could
	// pick for you, and guessing would route writes to a store nobody chose.
	CONF_DEFAULT_DATACONNECTION = "defaultdataconnection"
)

// NO_DATACONNECTION is the explicit declaration that a solution has no data layer at all. It is a
// value for CONF_DEFAULT_DATACONNECTION, and it lets such a solution boot without naming a
// connection it does not have -- while making the absence a STATEMENT rather than an omission.
//
// A solution declaring it may register NO data component: a component registering against a
// deployment that says it has no data layer is a contradiction, and failing at registration names
// the disagreement where it can be fixed, rather than at the first query.
const NO_DATACONNECTION = "<nodataconnection>"

func NotifyDelete(ctx core.RequestContext, objectType string, id string) {

}

/*
func GetFromCache(ctx core.RequestContext, objectType string, id string, object interface{}) bool {
	cachekey := components.GetCacheKey(objectType, id)
	return ctx.GetFromCache(cachekey, object)
}

func PutInCache(ctx core.RequestContext, objectType string, id string, object interface{}) {
	ctx.PutInCache(components.GetCacheKey(objectType, id), object)
}

func InvalidateCache(ctx core.RequestContext, objectType string, id string) {
	ctx.InvalidateCache(components.GetCacheKey(objectType, id))
}
*/
