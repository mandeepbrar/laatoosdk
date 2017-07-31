package data

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
)

/**
Base Component helps create a new data service
*/
type BaseComponent struct {
	core.Service
	Object                  string
	ObjectCreator           core.ObjectCreator
	ObjectCollectionCreator core.ObjectCollectionCreator
	ObjectConfig            *StorableConfig
	Auditable               bool
	SoftDelete              bool
	PreSave                 bool
	PostSave                bool
	PostLoad                bool
	PostUpdate              bool
	SoftDeleteField         string
	ObjectId                string
}

func (bc *BaseComponent) Initialize(ctx core.ServerContext) error {
	bc.SetComponent(ctx, true)
	bc.AddStringConfigurations(ctx, []string{CONF_DATA_OBJECT}, nil)
	bc.AddOptionalConfigurations(ctx, map[string]string{CONF_DATA_AUDITABLE: config.CONF_OBJECT_BOOL, CONF_DATA_POSTUPDATE: config.CONF_OBJECT_BOOL,
		CONF_DATA_POSTSAVE: config.CONF_OBJECT_BOOL, CONF_DATA_PRESAVE: config.CONF_OBJECT_BOOL, CONF_DATA_POSTLOAD: config.CONF_OBJECT_BOOL}, nil)
	return nil
}

func (bc *BaseComponent) Start(ctx core.ServerContext) error {
	object, _ := bc.GetConfiguration(ctx, CONF_DATA_OBJECT)
	bc.Object = object.(string)
	objectCreator, err := ctx.GetObjectCreator(bc.Object)
	if err != nil {
		return errors.BadArg(ctx, CONF_DATA_OBJECT, "Could not get Object creator for", object)
	}

	objectCollectionCreator, err := ctx.GetObjectCollectionCreator(bc.Object)
	if err != nil {
		return errors.BadArg(ctx, CONF_DATA_OBJECT, "Could not get Object collection creator for", object)
	}

	bc.ObjectCreator = objectCreator
	bc.ObjectCollectionCreator = objectCollectionCreator

	testObj := objectCreator()
	stor := testObj.(Storable)
	bc.ObjectConfig = stor.Config()

	bc.ObjectId = bc.ObjectConfig.IdField
	bc.SoftDeleteField = bc.ObjectConfig.SoftDeleteField

	if bc.SoftDeleteField == "" {
		bc.SoftDelete = false
	} else {
		bc.SoftDelete = true
	}

	auditable, ok := bc.GetConfiguration(ctx, CONF_DATA_AUDITABLE)
	if ok {
		bc.Auditable = auditable.(bool)
	} else {
		bc.Auditable = bc.ObjectConfig.Auditable
	}
	postsave, ok := bc.GetConfiguration(ctx, CONF_DATA_POSTSAVE)
	if ok {
		bc.PostSave = postsave.(bool)
	} else {
		bc.PostSave = bc.ObjectConfig.PostSave
	}
	postupdate, ok := bc.GetConfiguration(ctx, CONF_DATA_POSTUPDATE)
	if ok {
		bc.PostUpdate = postupdate.(bool)
	} else {
		bc.PostUpdate = bc.ObjectConfig.PostUpdate
	}
	presave, ok := bc.GetConfiguration(ctx, CONF_DATA_PRESAVE)
	if ok {
		bc.PreSave = presave.(bool)
	} else {
		bc.PreSave = bc.ObjectConfig.PreSave
	}
	postload, ok := bc.GetConfiguration(ctx, CONF_DATA_POSTLOAD)
	if ok {
		bc.PostLoad = postload.(bool)
	} else {
		bc.PostLoad = bc.ObjectConfig.PostLoad
	}

	return nil
}

func (bc *BaseComponent) GetDataServiceType() string {
	return ""
}

func (bc *BaseComponent) GetObject() string {
	return ""
}

func (bc *BaseComponent) GetCollection() string {
	return bc.ObjectConfig.Collection
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

func (bc *BaseComponent) CreateMulti(ctx core.RequestContext, items []Storable) error {
	return errors.NotImplemented(ctx, "CreateMulti")
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

func (bc *BaseComponent) CountGroups(ctx core.RequestContext, queryCond interface{}, groupids []string, group string) (res map[string]interface{}, err error) {
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
