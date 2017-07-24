package data

import (
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
	"laatoo/sdk/log"
)

/*
Plugins help a developer create layers of data services over one another
*/
type DataPlugin struct {
	*BaseComponent
	PluginDataComponent DataComponent
}

func NewDataPlugin(ctx core.ServerContext) *DataPlugin {
	return &DataPlugin{BaseComponent: &BaseComponent{}}
}
func NewDataPluginWithBase(ctx core.ServerContext, comp DataComponent) *DataPlugin {
	return &DataPlugin{BaseComponent: &BaseComponent{}, PluginDataComponent: comp}
}
func (svc *DataPlugin) Initialize(ctx core.ServerContext) error {
	err := svc.BaseComponent.Initialize(ctx)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	svc.AddStringConfigurations([]string{CONF_BASE_SVC}, []string{""})
	log.Error(ctx, "initialized ", "datacomponent", svc.PluginDataComponent)
	return nil
}

func (svc *DataPlugin) Start(ctx core.ServerContext) error {
	if svc.PluginDataComponent != nil {
		return nil
	}
	bsSvc, present := svc.GetConfiguration(CONF_BASE_SVC)
	if !present {
		return errors.MissingConf(ctx, CONF_BASE_SVC)
	}

	s, err := ctx.GetService(bsSvc.(string))
	if err != nil {
		return errors.BadConf(ctx, CONF_BASE_SVC)
	}

	PluginDataComponent, ok := s.(DataComponent)
	if !ok {
		return errors.BadConf(ctx, CONF_BASE_SVC)
	}
	svc.PluginDataComponent = PluginDataComponent
	return nil
}

func (svc *DataPlugin) GetCollection() string {
	return svc.PluginDataComponent.GetCollection()
}
func (svc *DataPlugin) CreateDBCollection(ctx core.RequestContext) error {
	return svc.PluginDataComponent.CreateDBCollection(ctx)
}

func (svc *DataPlugin) DropDBCollection(ctx core.RequestContext) error {
	return svc.PluginDataComponent.DropDBCollection(ctx)
}

func (svc *DataPlugin) DBCollectionExists(ctx core.RequestContext) (bool, error) {
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

//upsert an object ...insert if not there... update if there
func (svc *DataPlugin) UpsertId(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return svc.PluginDataComponent.UpsertId(ctx, id, newVals)
}

func (svc *DataPlugin) Update(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return svc.PluginDataComponent.Update(ctx, id, newVals)
}

func (svc *DataPlugin) Upsert(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return svc.PluginDataComponent.Upsert(ctx, queryCond, newVals)
}

func (svc *DataPlugin) UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return svc.PluginDataComponent.UpdateAll(ctx, queryCond, newVals)
}

//update objects by ids, fields to be updated should be provided as key value pairs
func (svc *DataPlugin) UpdateMulti(ctx core.RequestContext, ids []string, newVals map[string]interface{}) error {
	return svc.PluginDataComponent.UpdateMulti(ctx, ids, newVals)
}

//item must support Deleted field for soft deletes
func (svc *DataPlugin) Delete(ctx core.RequestContext, id string) error {
	return svc.PluginDataComponent.Delete(ctx, id)
}

//Delete object by ids
func (svc *DataPlugin) DeleteMulti(ctx core.RequestContext, ids []string) error {
	return svc.PluginDataComponent.DeleteMulti(ctx, ids)
}

//Delete object by condition
func (svc *DataPlugin) DeleteAll(ctx core.RequestContext, queryCond interface{}) ([]string, error) {
	return svc.PluginDataComponent.DeleteAll(ctx, queryCond)
}

func (svc *DataPlugin) GetById(ctx core.RequestContext, id string) (Storable, error) {
	return svc.PluginDataComponent.GetById(ctx, id)
}

//Get multiple objects by id
func (svc *DataPlugin) GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]Storable, error) {
	return svc.PluginDataComponent.GetMulti(ctx, ids, orderBy)
}

func (svc *DataPlugin) GetMultiHash(ctx core.RequestContext, ids []string) (map[string]Storable, error) {
	return svc.PluginDataComponent.GetMultiHash(ctx, ids)
}

func (svc *DataPlugin) Count(ctx core.RequestContext, queryCond interface{}) (count int, err error) {
	return svc.PluginDataComponent.Count(ctx, queryCond)
}

func (svc *DataPlugin) CountGroups(ctx core.RequestContext, queryCond interface{}, groupids []string, group string) (res map[string]interface{}, err error) {
	return svc.PluginDataComponent.CountGroups(ctx, queryCond, groupids, group)
}

func (svc *DataPlugin) GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.PluginDataComponent.GetList(ctx, pageSize, pageNum, mode, orderBy)
}

func (svc *DataPlugin) Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.PluginDataComponent.Get(ctx, queryCond, pageSize, pageNum, mode, orderBy)
}

//create condition for passing to data service
func (svc *DataPlugin) CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error) {
	return svc.PluginDataComponent.CreateCondition(ctx, operation, args...)
}
