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

	// CONF_DATA_NAMESPACE is the module/entity setting naming a component's namespace.
	CONF_DATA_NAMESPACE = "namespace"
)

// NAMESPACE_DEFAULT is the namespace a component occupies when nothing declares one, and the one
// a lookup means when it does not say.
//
// It is a real namespace rather than a derived value, and that is what makes existing behaviour
// survive. A draft derived the default from the component's dataconnection: in a single-connection
// deployment named "boltdb" every component would have landed in namespace "boltdb" and every
// lookup asking for the default would have found nothing -- 178 call sites failing in exactly the
// deployment shape that was supposed to see no change. Preserving behaviour by ABSENCE works;
// preserving it by a derivation that must be inverted at lookup does not.
const NAMESPACE_DEFAULT = "default"

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
