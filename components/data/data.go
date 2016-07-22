package data

import (
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
)

type Feature int

type ConditionType int

const (
	MATCHMULTIPLEVALUES ConditionType = iota // expects first value as field name and second value as array of values
	MATCHANCESTOR                            //expects collection name and id
	FIELDVALUE                               //expects map of field values
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
	GetObject() string
	//supported features
	Supports(Feature) bool
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
	GetMultiHash(ctx core.RequestContext, ids []string, orderBy string) (map[string]Storable, error)
	//Get multiple objects by id
	GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]Storable, error)
	//Get all object with given conditions
	Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Get a list of all items
	GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error)
}

type BaseComponent struct {
}

func (bc *BaseComponent) GetDataServiceType() string {
	return ""
}
func (bc *BaseComponent) GetObject() string {
	return ""
}

//supported features
func (bc *BaseComponent) Supports(Feature) bool {
	return false
}

//create condition for passing to data service
func (bc *BaseComponent) CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error) {
	return nil, errors.NotImplemented(ctx, "CreateCondition")
}

//save an object
func (bc *BaseComponent) Save(ctx core.RequestContext, item Storable) error {
	return errors.NotImplemented(ctx, "Save")
}

//adds an item to an array field
func (bc *BaseComponent) AddToArray(ctx core.RequestContext, id string, fieldName string, item interface{}) error {
	return errors.NotImplemented(ctx, "AddToArray")
}

//execute function
func (bc *BaseComponent) Execute(ctx core.RequestContext, name string, data interface{}, params map[string]interface{}) (interface{}, error) {
	return nil, errors.NotImplemented(ctx, "Execute")
}

//Store an object against an id
func (bc *BaseComponent) Put(ctx core.RequestContext, id string, item Storable) error {
	return errors.NotImplemented(ctx, "Put")
}

//Store multiple objects
func (bc *BaseComponent) PutMulti(ctx core.RequestContext, items []Storable) error {
	return errors.NotImplemented(ctx, "PutMulti")
}

//upsert an object ...insert if not there... update if there
func (bc *BaseComponent) UpsertId(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return errors.NotImplemented(ctx, "UpsertId")
}

//upsert an object ...insert if not there... update if there
func (bc *BaseComponent) Upsert(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return nil, errors.NotImplemented(ctx, "Upsert")
}

//update objects by ids, fields to be updated should be provided as key value pairs
func (bc *BaseComponent) UpdateMulti(ctx core.RequestContext, ids []string, newVals map[string]interface{}) error {
	return errors.NotImplemented(ctx, "UpdateMulti")
}

//update an object by ids, fields to be updated should be provided as key value pairs
func (bc *BaseComponent) Update(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return errors.NotImplemented(ctx, "Update")
}

//update with condition
func (bc *BaseComponent) UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return nil, errors.NotImplemented(ctx, "UpdateAll")
}

//Delete an object by id
func (bc *BaseComponent) Delete(ctx core.RequestContext, id string) error {
	return errors.NotImplemented(ctx, "Delete")
}

//Delete object by ids
func (bc *BaseComponent) DeleteMulti(ctx core.RequestContext, ids []string) error {
	return errors.NotImplemented(ctx, "DeleteMulti")
}

//delete with condition
func (bc *BaseComponent) DeleteAll(ctx core.RequestContext, queryCond interface{}) ([]string, error) {
	return nil, errors.NotImplemented(ctx, "DeleteAll")
}

//Get an object by id
func (bc *BaseComponent) GetById(ctx core.RequestContext, id string) (Storable, error) {
	return nil, errors.NotImplemented(ctx, "GetById")
}

//get storables in a hashtable
func (bc *BaseComponent) GetMultiHash(ctx core.RequestContext, ids []string, orderBy string) (map[string]Storable, error) {
	return nil, errors.NotImplemented(ctx, "GetMultiHash")
}

//Get multiple objects by id
func (bc *BaseComponent) GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]Storable, error) {
	return nil, errors.NotImplemented(ctx, "GetMulti")
}

//Get all object with given conditions
func (bc *BaseComponent) Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return nil, nil, -1, -1, errors.NotImplemented(ctx, "Get")
}

//Get a list of all items
func (bc *BaseComponent) GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return nil, nil, -1, -1, errors.NotImplemented(ctx, "GetList")

}
