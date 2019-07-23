package data

import "laatoo/sdk/server/core"

const (
	CONF_DATA_SVCS                = "dataservices"
	CONF_BASE_SVC                 = "baseservice"
	CONF_DATA_OBJECT              = "object"
	CONF_DATA_POSTSAVE            = "postsave"
	CONF_DATA_POSTLOAD            = "postload"
	CONF_DATA_MULTITENANT         = "multitenant"
	CONF_DATA_POSTUPDATE          = "postupdate"
	CONF_DATA_PRESAVE             = "presave"
	CONF_DATA_CACHEABLE           = "cacheable"
	CONF_DATA_AUDITABLE           = "auditable"
	CONF_DATA_REFOPS              = "refops"
	CONF_DATA_COLLECTION          = "collection"
	CONF_DATA_EMBEDDED_DOC_SEARCH = "embedded_doc_search"
	CONF_PRESAVE_MSG              = "storable_presave"
	CONF_PREUPDATE_MSG            = "storable_preupdate"
	CONF_POSTUPDATE_MSG           = "storable_postupdate"
	CONF_NEWOBJ_MSG               = "storable_new"
)

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
