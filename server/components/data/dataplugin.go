package data

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/utils"
)

/*
Plugins help a developer create layers of data services over one another
*/
type DataPlugin struct {
	core.Service
	PluginDataComponent DataComponent
}

func NewDataPlugin(ctx core.ServerContext) *DataPlugin {
	return &DataPlugin{}
}
func NewDataPluginWithBase(ctx core.ServerContext, comp DataComponent) *DataPlugin {
	return &DataPlugin{PluginDataComponent: comp}
}
func (svc *DataPlugin) Describe(ctx core.ServerContext) error {
	if svc.PluginDataComponent == nil {
		svc.AddStringConfiguration(ctx, CONF_BASE_SVC)
	}
	return nil
}

func (svc *DataPlugin) Initialize(ctx core.ServerContext, conf config.Config) error {
	if svc.PluginDataComponent != nil {
		return nil
	}

	bsSvc, _ := svc.GetStringConfiguration(ctx, CONF_BASE_SVC)
	s, err := ctx.GetService(bsSvc)
	if err != nil {
		return errors.BadConf(ctx, CONF_BASE_SVC)
	}

	dc, ok := s.(DataComponent)
	if !ok {
		return errors.BadConf(ctx, CONF_BASE_SVC)
	}
	svc.PluginDataComponent = dc

	return nil
}

func (svc *DataPlugin) GetObject() string {
	return svc.PluginDataComponent.GetObject()
}

func (svc *DataPlugin) GetCollection() string {
	return svc.PluginDataComponent.GetCollection()
}
func (svc *DataPlugin) CreateDBCollection(ctx core.ServerContext) error {
	return svc.PluginDataComponent.CreateDBCollection(ctx)
}

func (svc *DataPlugin) DropDBCollection(ctx core.ServerContext) error {
	return svc.PluginDataComponent.DropDBCollection(ctx)
}

func (svc *DataPlugin) DBCollectionExists(ctx core.ServerContext) (bool, error) {
	return svc.PluginDataComponent.DBCollectionExists(ctx)
}

func (svc *DataPlugin) GetDataServiceType() string {
	return svc.PluginDataComponent.GetDataServiceType()
}

func (svc *DataPlugin) Supports(feature Feature) bool {
	return svc.PluginDataComponent.Supports(feature)
}

func (svc *DataPlugin) Save(ctx core.RequestContext, item Storable) error {
	return svc.PluginDataComponent.Save(ctx, item)
}

func (svc *DataPlugin) PutMulti(ctx core.RequestContext, items []Storable) error {
	return svc.PluginDataComponent.PutMulti(ctx, items)
}

func (svc *DataPlugin) CreateMulti(ctx core.RequestContext, items []Storable) error {
	return svc.PluginDataComponent.CreateMulti(ctx, items)
}

func (svc *DataPlugin) Put(ctx core.RequestContext, id string, item Storable) error {
	return svc.PluginDataComponent.Put(ctx, id, item)
}

// upsert an object ...insert if not there... update if there
func (svc *DataPlugin) UpsertId(ctx core.RequestContext, id string, newVals utils.StringMap) error {
	return svc.PluginDataComponent.UpsertId(ctx, id, newVals)
}

func (svc *DataPlugin) Update(ctx core.RequestContext, id string, newVals utils.StringMap) error {
	return svc.PluginDataComponent.Update(ctx, id, newVals)
}

func (svc *DataPlugin) Upsert(ctx core.RequestContext, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error) {
	return svc.PluginDataComponent.Upsert(ctx, queryCond, newVals, getids)
}

func (svc *DataPlugin) UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals utils.StringMap, getids bool) ([]string, error) {
	return svc.PluginDataComponent.UpdateAll(ctx, queryCond, newVals, getids)
}

// update objects by ids, fields to be updated should be provided as key value pairs
func (svc *DataPlugin) UpdateMulti(ctx core.RequestContext, ids []string, newVals utils.StringMap) error {
	return svc.PluginDataComponent.UpdateMulti(ctx, ids, newVals)
}

// item must support Deleted field for soft deletes
func (svc *DataPlugin) Delete(ctx core.RequestContext, id string) error {
	return svc.PluginDataComponent.Delete(ctx, id)
}

// Delete object by ids
func (svc *DataPlugin) DeleteMulti(ctx core.RequestContext, ids []string) error {
	return svc.PluginDataComponent.DeleteMulti(ctx, ids)
}

// Delete object by condition
func (svc *DataPlugin) DeleteAll(ctx core.RequestContext, queryCond interface{}, getids bool) ([]string, error) {
	return svc.PluginDataComponent.DeleteAll(ctx, queryCond, getids)
}

func (svc *DataPlugin) GetById(ctx core.RequestContext, id string) (Storable, error) {
	return svc.PluginDataComponent.GetById(ctx, id)
}

// Get multiple objects by id
func (svc *DataPlugin) GetMulti(ctx core.RequestContext, ids []string, orderBy []string) ([]Storable, error) {
	return svc.PluginDataComponent.GetMulti(ctx, ids, orderBy)
}

// Gets the value of a key
func (svc *DataPlugin) GetValue(ctx core.RequestContext, key string) (interface{}, error) {
	return svc.PluginDataComponent.GetValue(ctx, key)
}

// Puts the value of a key
func (svc *DataPlugin) PutValue(ctx core.RequestContext, key string, value interface{}) error {
	return svc.PluginDataComponent.PutValue(ctx, key, value)
}

// Deletes the key
func (svc *DataPlugin) DeleteValue(ctx core.RequestContext, key string) error {
	return svc.PluginDataComponent.DeleteValue(ctx, key)
}

func (svc *DataPlugin) GetMultiHash(ctx core.RequestContext, ids []string) (map[string]Storable, error) {
	return svc.PluginDataComponent.GetMultiHash(ctx, ids)
}

func (svc *DataPlugin) Count(ctx core.RequestContext, queryCond interface{}) (count int, err error) {
	return svc.PluginDataComponent.Count(ctx, queryCond)
}

func (svc *DataPlugin) CountGroups(ctx core.RequestContext, queryCond interface{}, groupids []string, group string) (res utils.StringMap, err error) {
	return svc.PluginDataComponent.CountGroups(ctx, queryCond, groupids, group)
}

func (svc *DataPlugin) GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy []string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.PluginDataComponent.GetList(ctx, pageSize, pageNum, mode, orderBy)
}

func (svc *DataPlugin) Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy []string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.PluginDataComponent.Get(ctx, queryCond, pageSize, pageNum, mode, orderBy)
}

func (svc *DataPlugin) GetOne(ctx core.RequestContext, queryCond interface{}) (dataToReturn Storable, err error) {
	return svc.PluginDataComponent.GetOne(ctx, queryCond)
}

// create condition for passing to data service
func (svc *DataPlugin) CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error) {
	return svc.PluginDataComponent.CreateCondition(ctx, operation, args...)
}

func (svc *DataPlugin) AddToArray(ctx core.RequestContext, id string, fieldName string, item interface{}) error {
	return svc.PluginDataComponent.AddToArray(ctx, id, fieldName, item)
}

func (svc *DataPlugin) Execute(ctx core.RequestContext, name string, data interface{}, params utils.StringMap) (interface{}, error) {
	return svc.PluginDataComponent.Execute(ctx, name, data, params)
}
