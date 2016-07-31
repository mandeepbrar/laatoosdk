package data

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
)

const (
	CONF_DATA_OBJECT = "object"
)

type BaseComponent struct {
	Object                  string
	ObjectCreator           core.ObjectCreator
	ObjectCollectionCreator core.ObjectCollectionCreator
	ObjectConfig            *StorableConfig
}

func (bc *BaseComponent) Initialize(ctx core.ServerContext, conf config.Config) error {
	object, ok := conf.GetString(CONF_DATA_OBJECT)
	if !ok {
		return errors.MissingConf(ctx, CONF_DATA_OBJECT)
	}
	objectCreator, err := ctx.GetObjectCreator(object)
	if err != nil {
		return errors.BadArg(ctx, CONF_DATA_OBJECT, "Could not get Object creator for", object)
	}
	objectCollectionCreator, err := ctx.GetObjectCollectionCreator(object)
	if err != nil {
		return errors.BadArg(ctx, CONF_DATA_OBJECT, "Could not get Object collection creator for", object)
	}

	bc.Object = object
	bc.ObjectCreator = objectCreator
	bc.ObjectCollectionCreator = objectCollectionCreator

	testObj := objectCreator()
	stor := testObj.(Storable)
	bc.ObjectConfig = stor.Config()
	return nil
}

func (bc *BaseComponent) GetDataServiceType() string {
	return ""
}

func (bc *BaseComponent) GetObject() string {
	return ""
}

//get object creator
func (bc *BaseComponent) GetObjectCreator() core.ObjectCreator {
	return bc.ObjectCreator
}

//get object collection creator
func (bc *BaseComponent) GetObjectCollectionCreator() core.ObjectCollectionCreator {
	return bc.ObjectCollectionCreator
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
func (bc *BaseComponent) GetMultiHash(ctx core.RequestContext, ids []string) (map[string]Storable, error) {
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

func (bc *BaseComponent) Count(ctx core.RequestContext, queryCond interface{}) (count int, err error) {
	return -1, errors.NotImplemented(ctx, "Count")
}

func (bc *BaseComponent) CountGroups(ctx core.RequestContext, queryCond interface{}, group string) (res map[string]interface{}, err error) {
	return nil, errors.NotImplemented(ctx, "CountGroups")
}

func (bc *BaseComponent) CreateDBCollection(ctx core.RequestContext) error {
	return errors.NotImplemented(ctx, "CreateDBCollection")
}

func (bc *BaseComponent) DropDBCollection(ctx core.RequestContext) error {
	return errors.NotImplemented(ctx, "DropDBCollection")
}

func (bc *BaseComponent) DBCollectionExists(ctx core.RequestContext) (bool, error) {
	return false, errors.NotImplemented(ctx, "DBCollectionExists")
}
