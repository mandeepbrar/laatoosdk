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
	//create condition for passing to data service
	CreateCondition(ctx core.RequestContext, obj string, operation data.ConditionType, args ...interface{}) (interface{}, error)

	Save(ctx core.RequestContext, obj string, item data.Storable) error
	//Store an object against an id
	Put(ctx core.RequestContext, obj string, id string, item data.Storable) error
	//Store multiple objects
	CreateMulti(ctx core.RequestContext, obj string, items []data.Storable) error
	//Store multiple objects
	PutMulti(ctx core.RequestContext, obj string, items []data.Storable) error
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
	GetById(ctx core.RequestContext, obj string, id string) (data.Storable, error)
	//get storables in a hashtable
	GetMultiHash(ctx core.RequestContext, obj string, ids []string) (map[string]data.Storable, error)
	//Get multiple objects by id
	GetMulti(ctx core.RequestContext, obj string, ids []string, orderBy interface{}) ([]data.Storable, error)
	//Gets the value of a key.
	GetValue(ctx core.RequestContext, obj string, key string) (interface{}, error)
	//Puts the value of a key
	PutValue(ctx core.RequestContext, obj string, key string, value interface{}) error
	//Deletes the key
	DeleteValue(ctx core.RequestContext, obj string, key string) error

	FetchDataset(ctx core.RequestContext, dsname string, params utils.StringsMap, pageSize int, pageNum int) (dataToReturn []data.Storable, ids []string, totalrecs int, recsreturned int, err error)

	//Count all object with given condition
	Count(ctx core.RequestContext, obj string, queryCond interface{}) (count int, err error)
	CountGroups(ctx core.RequestContext, obj string, queryCond interface{}, groupids []string, group string) (res utils.StringMap, err error)

	Transaction(ctx core.RequestContext, obj string, callback func(ctx core.RequestContext) error) error

	//Get all object with given conditions
	Get(ctx core.RequestContext, obj string, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy interface{}) (dataToReturn []data.Storable, ids []string, totalrecs int, recsreturned int, err error)
	//Get one record satisfying condition
	GetOne(ctx core.RequestContext, obj string, queryCond interface{}) (dataToReturn data.Storable, err error)
	//Get a list of all items
	GetList(ctx core.RequestContext, obj string, pageSize int, pageNum int, mode string, orderBy interface{}) (dataToReturn []data.Storable, ids []string, totalrecs int, recsreturned int, err error)
}
