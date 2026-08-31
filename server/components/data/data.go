package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Feature int

type SortType string

const (
	SORTASC  SortType = "asc"
	SORTDESC          = "desc"
)

const (
	InQueries Feature = iota
	Ancestors
	EmbeddedSearch
)

const (
	DATASERVICE_TYPE_NOSQL = "SERVICE_TYPE_NOSQL"
	DATASERVICE_TYPE_SQL   = "SERVICE_TYPE_SQL"
	DATASERVICE_TYPE_KV    = "SERVICE_TYPE_KV"
	DATASERVICE_TYPE_CACHE = "SERVICE_TYPE_CACHE"
	DATA_PAGENUM           = "pagenum"
	DATA_PAGESIZE          = "pagesize"
	DATA_RECSRETURNED      = "records"
	DATA_TOTALRECS         = "totalrecords"
)

type DataEventType string

// Standard Data Events
const (
	EventDataCreated DataEventType = "data.object.created"
	EventDataUpdated DataEventType = "data.object.updated"
	EventDataDeleted DataEventType = "data.object.deleted"
)

type Dataset struct {
	Name       string
	Properties utils.StringsMap
	Entity     string
	// Connection names the dataconnection holding Entity, the other half of the key the component
	// registry is built on. Declared in the dataset YAML as `Connection`, and empty means the
	// entity is registered on exactly one connection and resolves there.
	Connection string
	// ConnectionNative reports whether this dataset is PROVABLY confined to one dataconnection.
	//
	// Computed at startup by the DataManager, not declared: it is a fact about where components
	// registered, which a dataset author cannot know and must not have to restate.
	//
	// It is CONSERVATIVE, and the asymmetry is deliberate. True means proven: the dataset reaches
	// exactly one connection and nothing it does can cross a store. False means NOT PROVEN, which
	// includes datasets that will in fact run on one connection.
	//
	// The reason it cannot be exact is structural: Expansion.Field names a FIELD, and the entity
	// that field reaches is carried on the StorableRef of a fetched record, not in any
	// declaration. core.StorableConfig holds no reference metadata, so at startup there is
	// nothing to resolve a hop target against. A dataset with hops is therefore unprovable here
	// even when every hop stays home.
	//
	// Read it as a guarantee when true and as "ask at runtime" when false. Treating false as
	// "spans stores" would be wrong.
	ConnectionNative bool
	QueryType  string
	QueryData  interface{}
	Params     utils.StringsMap
	Cache      bool
	Dao        string
	Permission string
	Module     core.Module
}

type VectorResult struct {
	Item  core.Storable
	Score float64
	Dist  float64
}

// Service that provides data from various data sources
// Service interface that needs to be implemented by any data service
type DataComponent interface {
	core.Service

	//GetConnectionName reports the dataconnection this component is bound to, forming the other
	//half of the (dataconnection, object) key the DataManager registry uses.
	//
	//The value is the FACTORY INSTANCE name, which is what a dataconnection is: a provider's
	//factory config sets `name:` to `{{if exists "dataconnection"}} {{var "dataconnection"}}
	//{{else}} <plugin> {{end}}`, so the factory registers under the deployment's dataconnection
	//value and falls back to the plugin name when the deployment declares none. A component is
	//always created by a factory, so this is never empty and there is no default to substitute.
	//
	//IT IS NOT core.Service.GetNamespace, which this interface also carries through the embedded
	//Service above. That one is the SECURITY namespace: the server reads it on every request and
	//hands it to SecurityHandler.AuthorizeService (laatooserver core/serviceelement.go:203).
	//Answering the storage question with it would silently relocate authorization.
	//
	//Required rather than optional: an optional interface leaves a provider that forgets it
	//unregistered or misrouted, which is the "both forms work and nobody notices" shape the
	//platform avoids elsewhere. Every implementor in the platform embeds
	//datacommon.BaseComponent, so one implementation there answers for all seven providers and
	//the masterdata wrapper.
	GetConnectionName() string

	//IsMyDataConnection reports whether this component can reach the named dataconnection.
	//
	//It is NOT derivable from GetConnectionName, and that is the whole reason it exists.
	//GetConnectionName answers "where am I" with a single string; this answers "can I reach
	//there", and the two differ whenever one provider spans several logical connections in one
	//physical store -- two schemas in one Postgres, two buckets in one Couchbase. Such a provider
	//is a single store for join purposes and should say so.
	//
	//THE CALLER THIS EXISTS FOR IS JOIN FEASIBILITY. Deciding whether a reference or a traversal
	//can compile natively is a question about whether one provider can reach both sides, and only
	//the provider knows. The registry cannot answer it: it knows which component holds each
	//entity, not which of them share a store.
	//
	//BaseComponent's default is strict equality, which is correct for every provider that maps one
	//connection to one store -- all of them today. A provider widens it only when it genuinely
	//spans more, and a provider that widens it dishonestly produces a join that compiles and reads
	//from the wrong store, so the same honesty rule applies here as to query capabilities.
	IsMyDataConnection(ctx core.ServerContext, connection string) bool

	//GetDataServiceType reports which storage family backs this component, as one of the
	//DATASERVICE_TYPE_* constants above. Providers return a fixed value — mongodatabase returns
	//DATASERVICE_TYPE_NOSQL (mongodataservice.go:131-133) — and it is a static property of the
	//plugin, not of the entity or the connection.
	//
	//BaseComponent's default returns the EMPTY STRING, not one of the constants
	//(datastore/dev/plugins/common/basecomponent.go:127-129), so a provider that forgets to
	//override it registers successfully and reports a type matching no constant. The only place
	//the value is consumed is the registration log line
	//(laatooserver/src/core/datamanager.go:107), which is why nothing catches that.
	GetDataServiceType() string
	//object on which service operates
	GetObject() string
	//collection for the service
	GetCollection() string
	//create object
	CreateObject(ctx core.RequestContext) interface{}
	//create object collection
	CreateObjectCollection(ctx core.RequestContext, len int) interface{}
	//create object pointers collection using factory
	CreateObjectPointersCollection(ctx core.RequestContext, len int) interface{}
	//object factory for the data object
	GetObjectFactory() core.ObjectFactory
	//supported features
	Supports(Feature) bool
	//creates a collection
	CreateDBCollection(ctx core.ServerContext) error
	//drops a collection
	DropDBCollection(ctx core.ServerContext) error
	//collection exists
	DBCollectionExists(ctx core.ServerContext) (bool, error)
	//create condition from field/value pairs, combined with equality and and. This is the
	//shorthand the large majority of callers use, and it stays the shortest way to express the
	//common case: build the map, omit the entries you do not want. It is sugar over
	//CreateQueryCondition, not a separate mechanism.
	CreateCondition(ctx core.RequestContext, args utils.StringMap) (interface{}, error)
	//compiles and binds a query in one step, for shapes built per request that the shorthand
	//cannot express — ranges, disjunction, substring matching, nesting.
	CreateQueryCondition(ctx core.RequestContext, query *Query, params utils.StringsMap) (interface{}, error)
	//compiles a query into a reusable provider-native form. Called once, when the query's shape
	//is known; implementations compile each predicate to an independently composable fragment so
	//that binding can assemble only the predicates that survive parameter resolution.
	CompileQuery(ctx core.ServerContext, query *Query) (interface{}, error)
	//binds parameters to a previously compiled query, producing a condition for the query
	//methods below. Optional predicates whose parameters are absent are dropped here.
	BindQuery(ctx core.RequestContext, compiled interface{}, params utils.StringsMap) (interface{}, error)
	//reports whether this provider can compile a given query capability natively. A provider
	//must not silently reduce a capability it lacks to a weaker one — it declares it here and
	//rejects it at compile time.
	SupportsQuery(capability QueryCapability) bool
	//starts a chained query against this component: build with Where/Through/Expanding, then end
	//with All, One, Count or Condition. It is the convenience path for a query built per request;
	//a fixed-shape query run many times still belongs in CompileQuery once and BindQuery per
	//request, because a builder compiles on every terminal.
	//
	//Every provider inherits BaseComponent's implementation, which binds the builder to the
	//concrete component — no provider implements this itself.
	CreateQuery(ctx core.RequestContext) *QueryBuilder
	//
	//A caller holding query TEXT rather than predicates uses elements.DataManager.CreateTextQuery
	//instead. Text does NOT belong on this interface: reading it means resolving a QueryComponent
	//by name, the registry of those lives on the DataManager, and a component reaching back up to
	//its own manager to answer a call made on itself is a cycle with no purpose. This is why the
	//CreateODataQuery that stood here was declared and never implemented — every implementor
	//failed the chain with Core_Not_Implemented, because a component genuinely cannot parse.
	//save an object
	Save(ctx core.RequestContext, item core.Storable) error
	//AddToArray appends an item to a list-valued field of one record, without loading it.
	//
	//NO PROVIDER IMPLEMENTS IT. Every shipped datastore either inherits BaseComponent's
	//errors.NotImplemented (basecomponent.go:412-414) — mongodatabase, sqldatabase, jsonbdatabase,
	//boltdatabase and couchbasedatabase all do — or overrides it with the same refusal
	//(gaedatastore gaedataservice.go:674-676, gaefirestore gaefireservice.go:737-739). Calling it
	//returns a Core_Not_Implemented error at runtime, never nil. To append to a list today, read
	//the record, append in Go, and Put it back.
	AddToArray(ctx core.RequestContext, id string, fieldName string, item interface{}) error
	//Execute runs a provider-native named operation — a stored procedure, a server-side script —
	//that the portable API cannot express.
	//
	//NO PROVIDER IMPLEMENTS IT either; the situation is exactly AddToArray's
	//(basecomponent.go:417-419, gaedataservice.go:678-680, gaefireservice.go:742-744). It returns
	//(nil, Core_Not_Implemented) on every shipped store. Because the first return is interface{},
	//a caller that ignores the error gets a nil it will dereference.
	Execute(ctx core.RequestContext, name string, data interface{}, params utils.StringMap) (interface{}, error)
	//Store an object against an id
	Put(ctx core.RequestContext, id string, item core.Storable) error
	//Store multiple objects
	CreateMulti(ctx core.RequestContext, items []core.Storable) error
	//Store multiple objects
	PutMulti(ctx core.RequestContext, items []core.Storable) error
	//upsert an object by id, fields to be updated should be provided as key value pairs
	UpsertId(ctx core.RequestContext, id string, newVals utils.StringMap) error
	//upsert by condition
	Upsert(ctx core.RequestContext, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error)
	//update objects by ids, fields to be updated should be provided as key value pairs
	UpdateMulti(ctx core.RequestContext, ids []string, newVals utils.StringMap) error
	//update an object by ids, fields to be updated should be provided as key value pairs
	Update(ctx core.RequestContext, id string, newVals utils.StringMap) error
	//update with condition
	UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error)
	//Delete an object by id
	Delete(ctx core.RequestContext, id string) error
	//Delete object by ids
	DeleteMulti(ctx core.RequestContext, ids []string) error
	//delete with condition
	DeleteAll(ctx core.RequestContext, queryCond interface{}, getids bool) ([]string, error)
	//Get an object by id
	GetById(ctx core.RequestContext, id string, dao string) (core.Storable, error)
	//get storables in a hashtable
	GetMultiHash(ctx core.RequestContext, props []string, ids []string, dao string) (map[string]core.Storable, error)
	//Get multiple objects by id
	GetMulti(ctx core.RequestContext, props []string, ids []string, orderBy []string, dao string) ([]core.Storable, error)
	//Gets the value of a key.
	GetValue(ctx core.RequestContext, key string) (interface{}, error)
	//Puts the value of a key
	PutValue(ctx core.RequestContext, key string, value interface{}) error
	//Deletes the key
	DeleteValue(ctx core.RequestContext, key string) error

	//Count all object with given condition
	Count(ctx core.RequestContext, queryCond interface{}) (count int, err error)
	CountGroups(ctx core.RequestContext, queryCond interface{}, groupids []string, group string) (res utils.StringMap, err error)

	Transaction(ctx core.RequestContext, callback func(ctx core.RequestContext) error) error

	//Get all object with given conditions
	Get(ctx core.RequestContext, props []string, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy []string, dao string) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Get one record satisfying condition
	GetOne(ctx core.RequestContext, props []string, queryCond interface{}, dao string) (dataToReturn core.Storable, err error)
	//Get a list of all items
	GetList(ctx core.RequestContext, props []string, pageSize int, pageNum int, mode string, orderBy []string, dao string) (dataToReturn []core.Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Vector Search
	VectorSearch(ctx core.RequestContext, vector []float32, limit int, filter interface{}) ([]VectorResult, error)
	//Subscribe to data events
	Subscribe(ctx core.RequestContext, obj string, eventType DataEventType, handler core.MessageListener) error
}
