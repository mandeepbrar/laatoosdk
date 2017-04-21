package data

import "laatoo/sdk/core"

type Feature int

type ConditionType int

const (
	MATCHMULTIPLEVALUES ConditionType = iota // expects first value as field name and second value as array of values
	MATCHANCESTOR                            //expects collection name and id
	FIELDVALUE                               //expects map of field values
	COMBINECONDTITIONS                       //combine conditions
)

const (
	InQueries Feature = iota
	Ancestors
)

const (
	DATASERVICE_TYPE_NOSQL = "SERVICE_TYPE_NOSQL"
	DATASERVICE_TYPE_SQL   = "SERVICE_TYPE_SQL"
	DATASERVICE_TYPE_CACHE = "SERVICE_TYPE_CACHE"
	DATA_PAGENUM           = "pagenum"
	DATA_PAGESIZE          = "pagesize"
	DATA_RECSRETURNED      = "records"
	DATA_TOTALRECS         = "totalrecords"
)

//Service that provides data from various data sources
//Service interface that needs to be implemented by any data service
type DataComponent interface {
	GetDataServiceType() string
	//object on which service operates
	GetObject() string
	//collection for the service
	GetCollection() string
	//get object creator
	GetObjectCreator() core.ObjectCreator
	//get object collection creator
	GetObjectCollectionCreator() core.ObjectCollectionCreator
	//supported features
	Supports(Feature) bool
	//creates a collection
	CreateDBCollection(ctx core.RequestContext) error
	//drops a collection
	DropDBCollection(ctx core.RequestContext) error
	//collection exists
	DBCollectionExists(ctx core.RequestContext) (bool, error)
	//create condition for passing to data service
	CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error)
	//save an object
	Save(ctx core.RequestContext, item Storable) error
	//adds an item to an array field
	AddToArray(ctx core.RequestContext, id string, fieldName string, item interface{}) error
	//execute function
	Execute(ctx core.RequestContext, name string, data interface{}, params map[string]interface{}) (interface{}, error)
	//Store an object against an id
	Put(ctx core.RequestContext, id string, item Storable) error
	//Store multiple objects
	CreateMulti(ctx core.RequestContext, items []Storable) error
	//Store multiple objects
	PutMulti(ctx core.RequestContext, items []Storable) error
	//upsert an object by id, fields to be updated should be provided as key value pairs
	UpsertId(ctx core.RequestContext, id string, newVals map[string]interface{}) error
	//upsert by condition
	Upsert(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error)
	//update objects by ids, fields to be updated should be provided as key value pairs
	UpdateMulti(ctx core.RequestContext, ids []string, newVals map[string]interface{}) error
	//update an object by ids, fields to be updated should be provided as key value pairs
	Update(ctx core.RequestContext, id string, newVals map[string]interface{}) error
	//update with condition
	UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error)
	//Delete an object by id
	Delete(ctx core.RequestContext, id string) error
	//Delete object by ids
	DeleteMulti(ctx core.RequestContext, ids []string) error
	//delete with condition
	DeleteAll(ctx core.RequestContext, queryCond interface{}) ([]string, error)
	//Get an object by id
	GetById(ctx core.RequestContext, id string) (Storable, error)
	//get storables in a hashtable
	GetMultiHash(ctx core.RequestContext, ids []string) (map[string]Storable, error)
	//Get multiple objects by id
	GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]Storable, error)

	//Count all object with given condition
	Count(ctx core.RequestContext, queryCond interface{}) (count int, err error)
	CountGroups(ctx core.RequestContext, queryCond interface{}, groupids []string, group string) (res map[string]interface{}, err error)

	//Get all object with given conditions
	Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Get a list of all items
	GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error)
}
