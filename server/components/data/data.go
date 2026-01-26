package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Feature int

type ConditionType int

const (
	ODATA ConditionType = iota // expects first value as field name and second value as array of values
	MATCHMULTIPLEVALUES
	MATCHANCESTOR //expects collection name and id
	FIELDVALUE    //expects map of field values
)

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

type Dataset struct {
	Name       string
	Properties utils.StringsMap
	Entity     string
	QueryType  string
	QueryData  interface{}
	Params     utils.StringsMap
	Cache      bool
	Dao        string
	Permission string
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
	//create condition for passing to data service
	CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error)
	//save an object
	Save(ctx core.RequestContext, item core.Storable) error
	//adds an item to an array field
	AddToArray(ctx core.RequestContext, id string, fieldName string, item interface{}) error
	//execute function
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
}

