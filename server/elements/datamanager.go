package elements

import (
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type DataManager interface {
	core.ServerElement

	//register the data component for an object type
	RegisterDataComponent(ctx core.ServerContext, obj string, comp data.DataComponent) error
	//get component registered for an entity
	GetRegisteredComponent(ctx core.ServerContext, obj string) (data.DataComponent, error)
	//create condition from field/value pairs combined with equality — the shorthand, unchanged
	//in shape from what callers have always written
	CreateCondition(ctx core.RequestContext, obj string, args utils.StringMap) (interface{}, error)
	//create condition from a query, for shapes the shorthand cannot express. Callers needing
	//compile-once/bind-per-request reach the component through GetRegisteredComponent.
	CreateQueryCondition(ctx core.RequestContext, obj string, query *data.Query, params utils.StringsMap) (interface{}, error)
	//start a chained query against the component registered for obj — the DataManager-scoped
	//twin of DataComponent.CreateQuery, taking the entity name a caller here does not otherwise
	//have a component for. Build with Where/Through/Expanding, end with All, One, Count or
	//Condition.
	//
	//An obj with no registered component yields a builder that fails at its terminal, not a nil,
	//so a chain reports the missing component where the caller is already checking an error.
	CreateQuery(ctx core.RequestContext, obj string) *data.QueryBuilder
	//start a chained query from OData filter text against the component registered for obj. The
	//parameter is a string and this contract carries no parser; parsing belongs to the component.
	CreateODataQuery(ctx core.RequestContext, obj string, odataQuery string) *data.QueryBuilder

	//Save writes an item through the data component registered for obj, after running the
	//component's configured hooks: presave message + Storable.PreSave, tenant stamping when the
	//component is multitenant, audit stamping when it is trackable, then the write, then
	//Storable.PostSave and the new-object message. See mongodatabase mongodataservice.go:151-201
	//for the canonical order; every provider follows it.
	//
	//SAVE IS NOT A UNIVERSAL UPSERT, whatever the older prose says. The underlying write differs
	//by provider: mongodatabase issues InsertOne (mongodataservice.go:174), sqldatabase issues
	//gorm Create (sqldataservice.go:252) and boltdatabase issues Insert — all three FAIL on an id
	//that already exists — while couchbasedatabase issues Upsert (couchbasedataservice.go:164) and
	//silently replaces. Use Put when you mean create-or-replace: it is upsert on every provider.
	//
	//It resolves the component by name, so an obj with no registered data component returns a
	//Core_Not_Found error rather than doing nothing (laatooserver/src/core/datamanager.go:111-117).
	Save(ctx core.RequestContext, obj string, item core.Storable) error
	//Store an object against an id
	Put(ctx core.RequestContext, obj string, id string, item core.Storable) error
	//Store multiple objects
	CreateMulti(ctx core.RequestContext, obj string, items []core.Storable) error
	//Store multiple objects
	PutMulti(ctx core.RequestContext, obj string, items []core.Storable) error
	//upsert an object by id, fields to be updated should be provided as key value pairs
	UpsertId(ctx core.RequestContext, obj string, id string, newVals utils.StringMap) error
	//upsert by condition
	Upsert(ctx core.RequestContext, obj string, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error)
	//update objects by ids, fields to be updated should be provided as key value pairs
	UpdateMulti(ctx core.RequestContext, obj string, ids []string, newVals utils.StringMap) error
	//update an object by ids, fields to be updated should be provided as key value pairs
	Update(ctx core.RequestContext, obj string, id string, newVals utils.StringMap) error
	//update with condition
	UpdateAll(ctx core.RequestContext, obj string, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error)
	//Delete an object by id
	Delete(ctx core.RequestContext, obj string, id string) error
	//Delete object by ids
	DeleteMulti(ctx core.RequestContext, obj string, ids []string) error
	//delete with condition
	DeleteAll(ctx core.RequestContext, obj string, queryCond interface{}, getids bool) ([]string, error)
	//Get an object by id
	GetById(ctx core.RequestContext, obj string, id string, dao string) (core.Storable, error)
	//get storables in a hashtable
	GetMultiHash(ctx core.RequestContext, props []string, obj string, ids []string, dao string) (map[string]core.Storable, error)
	//Get multiple objects by id
	GetMulti(ctx core.RequestContext, props []string, obj string, ids []string, orderBy []string, dao string) ([]core.Storable, error)
	//Gets the value of a key.
	GetValue(ctx core.RequestContext, obj string, key string) (interface{}, error)
	//Puts the value of a key
	PutValue(ctx core.RequestContext, obj string, key string, value interface{}) error
	//Deletes the key
	DeleteValue(ctx core.RequestContext, obj string, key string) error

	//FetchDataset runs a dataset declared in a plugin's registry/datasets and returns its page.
	//
	//It is the only entry point here that enforces a permission of its own: a dataset declaring
	//Permission is checked with ctx.HasPermission before anything runs, and an unauthorised caller
	//gets Core_Error_Unauthorized (laatooserver/src/common/fieldvalueds.go:104-106). Nothing else
	//on this interface consults permissions.
	//
	//The dataset's query is compiled ONCE, on first fetch rather than at load — datasets are read
	//before data components are registered, so nothing exists to compile against at load time
	//(fieldvalueds.go:122-128). A compile error is therefore first seen by the first caller, and
	//is then returned to EVERY later caller because sync.Once does not retry.
	//
	//params binds the dataset's declared Params. A filter not marked Optional whose parameter is
	//absent fails with Core_Missing_Arg (fieldvalueds.go:114-120) instead of matching nothing;
	//an Optional filter with an absent parameter is dropped from the query, which WIDENS the
	//result set.
	//
	//pageNum is 1-BASED. Providers compute skip as (pageNum-1)*pageSize
	//(mongodatabase mongodataservice_get.go:220), so 0 yields a negative skip that the store
	//rejects. Pass -1, -1 for everything.
	//
	//An unknown dsname returns Core_Not_Found (laatooserver/src/core/datamanager.go:99).
	FetchDataset(ctx core.RequestContext, dsname string, params utils.StringsMap, pageSize int, pageNum int) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)

	//Count all object with given condition
	Count(ctx core.RequestContext, obj string, queryCond interface{}) (count int, err error)
	//CountGroups returns per-group counts for the records matching queryCond: a map keyed by the
	//distinct values of the group field, each holding that group's count. groupids, when non-nil,
	//restricts the result to those group values.
	//
	//IT IS NOT SUPPORTED ON EVERY STORE. sqldatabase (sqldataservice_get.go:165), jsonbdatabase
	//and boltdatabase implement it; mongodatabase and couchbasedatabase implement nothing and so
	//inherit BaseComponent's errors.NotImplemented (basecomponent.go:520-522), and gaedatastore
	//(gaedataservice.go:682-684) and gaefirestore override it with the same refusal. On those five
	//it returns (nil, Core_Not_Implemented) — a nil map, so a caller that ignores the error reads
	//zero for every group rather than failing.
	CountGroups(ctx core.RequestContext, obj string, queryCond interface{}, groupids []string, group string) (res utils.StringMap, err error)

	//Transaction runs callback inside a transaction opened on the data component registered for
	//obj, committing when callback returns nil and rolling back when it returns an error.
	//
	//THE CALLBACK IS HANDED A DIFFERENT CONTEXT AND MUST USE IT. Providers bind the transaction to
	//that context and nowhere else — sqldatabase stores the gorm tx as a context value
	//(sqldataservice.go:680-682), mongodatabase wraps it in a driver SessionContext
	//(mongodataservice.go:678-706). Work issued through the OUTER ctx inside the callback runs
	//outside the transaction and is committed independently; nothing detects or reports that, so a
	//rollback silently leaves it behind.
	//
	//Scope is one connection, not the whole data layer: entities served by a different component
	//or a different connection do not join. And it is not reentrant on boltdatabase, where bbolt
	//permits a single writer — a nested write opens a second transaction and self-deadlocks with
	//no error and no timeout (see the note on kvDataService.insertInTx).
	//
	//sqldatabase discards the result of Rollback (sqldataservice.go:685): the callback's error is
	//what surfaces, and a rollback that itself failed is invisible.
	Transaction(ctx core.RequestContext, obj string, callback func(ctx core.RequestContext) error) error

	//Get all object with given conditions
	Get(ctx core.RequestContext, props []string, obj string, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy []string, dao string) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Get one record satisfying condition
	GetOne(ctx core.RequestContext, props []string, obj string, queryCond interface{}, dao string) (dataToReturn core.Storable, err error)
	//Get a list of all items
	GetList(ctx core.RequestContext, props []string, obj string, pageSize int, pageNum int, mode string, orderBy []string, dao string) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)

	//Vector Search
	VectorSearch(ctx core.RequestContext, obj string, vector []float32, limit int, filter interface{}) ([]data.VectorResult, error)
	//Subscribe to data events
	Subscribe(ctx core.RequestContext, obj string, eventType data.DataEventType, handler core.MessageListener) error

	// ---- Namespace-carrying variants -------------------------------------------------------
	//
	// The registry is keyed on (namespace, object), and every method above means the DEFAULT
	// namespace. These say which one instead.
	//
	// They are PARALLEL METHODS rather than widened signatures, and Go has no overloading, so
	// they carry their own names. The reason is measured: 178 production call sites across ~65
	// files use the methods above, none of them generated, and widening would have edited every
	// one to insert a constant. The methods above are not a degraded form -- for the
	// single-namespace deployment, which is most of them, they are exactly right.
	//
	// Only the methods measured as actually called have a variant. The unused tail of the
	// interface stays as it is until something needs it; adding a variant nobody calls is
	// surface with no reader.
	//
	// A namespace of "" is read as NAMESPACE_DEFAULT, so these are safe to call unconditionally
	// from code that may or may not have a namespace in hand.
	CreateConditionInNamespace(ctx core.RequestContext, namespace string, obj string, args utils.StringMap) (interface{}, error)
	CreateQueryConditionInNamespace(ctx core.RequestContext, namespace string, obj string, query *data.Query, params utils.StringsMap) (interface{}, error)
	SaveInNamespace(ctx core.RequestContext, namespace string, obj string, item core.Storable) error
	PutInNamespace(ctx core.RequestContext, namespace string, obj string, id string, item core.Storable) error
	UpsertIdInNamespace(ctx core.RequestContext, namespace string, obj string, id string, newVals utils.StringMap) error
	UpdateInNamespace(ctx core.RequestContext, namespace string, obj string, id string, newVals utils.StringMap) error
	DeleteInNamespace(ctx core.RequestContext, namespace string, obj string, id string) error
	GetByIdInNamespace(ctx core.RequestContext, namespace string, obj string, id string, dao string) (core.Storable, error)
	GetInNamespace(ctx core.RequestContext, namespace string, props []string, obj string, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy []string, dao string) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)
	GetOneInNamespace(ctx core.RequestContext, namespace string, props []string, obj string, queryCond interface{}, dao string) (dataToReturn core.Storable, err error)

	// RegisterDataComponentInNamespace registers a component under an explicit namespace.
	//
	// RegisterDataComponent above places the component in the namespace it declares through
	// data.NamespacedComponent, or NAMESPACE_DEFAULT when it declares none. This is for a caller
	// that knows better than the component does.
	//
	// EITHER WAY, A DUPLICATE (namespace, object) IS REFUSED rather than replacing what is
	// there. Until this change the registry did an unconditional Store, so two components
	// claiming one name both logged success and every write went to whichever registered last.
	RegisterDataComponentInNamespace(ctx core.ServerContext, namespace string, obj string, comp data.DataComponent) error
	// GetRegisteredComponentInNamespace resolves a component in an explicit namespace.
	GetRegisteredComponentInNamespace(ctx core.ServerContext, namespace string, obj string) (data.DataComponent, error)
}
