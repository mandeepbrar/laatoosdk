package data

import (
	"laatoo/framework/services/data/common"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
)

type DataPlugin struct {
	*BaseComponent
	dataServiceName string
	DataComponent   DataComponent
}

func NewDataPlugin(ctx core.ServerContext) *DataPlugin {
	return &DataPlugin{BaseComponent: &BaseComponent{}}
}

func (svc *DataPlugin) Initialize(ctx core.ServerContext, conf config.Config) error {
	err := svc.BaseComponent.Initialize(ctx, conf)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	bsSvc, ok := conf.GetString(common.CONF_BASE_SVC)
	if !ok {
		return errors.MissingConf(ctx, common.CONF_BASE_SVC)
	}
	svc.dataServiceName = bsSvc
	return nil
}

func (svc *DataPlugin) Start(ctx core.ServerContext) error {
	s, err := ctx.GetService(svc.dataServiceName)
	if err != nil {
		return errors.BadConf(ctx, common.CONF_BASE_SVC)
	}
	DataComponent, ok := s.(DataComponent)
	if !ok {
		return errors.BadConf(ctx, common.CONF_BASE_SVC)
	}
	svc.DataComponent = DataComponent
	return nil
}

func (svc *DataPlugin) Invoke(ctx core.RequestContext) error {
	return nil
}

func (svc *DataPlugin) CreateDBCollection(ctx core.RequestContext) error {
	return svc.DataComponent.CreateDBCollection(ctx)
}

func (svc *DataPlugin) DropDBCollection(ctx core.RequestContext) error {
	return svc.DataComponent.DropDBCollection(ctx)
}

func (svc *DataPlugin) DBCollectionExists(ctx core.RequestContext) (bool, error) {
	return svc.DataComponent.DBCollectionExists(ctx)
}

func (svc *DataPlugin) GetDataServiceType() string {
	return svc.DataComponent.GetDataServiceType()
}

func (svc *DataPlugin) Supports(feature Feature) bool {
	return svc.DataComponent.Supports(feature)
}

func (svc *DataPlugin) Save(ctx core.RequestContext, item Storable) error {
	return svc.DataComponent.Save(ctx, item)
}

func (svc *DataPlugin) PutMulti(ctx core.RequestContext, items []Storable) error {
	return svc.DataComponent.PutMulti(ctx, items)
}

func (svc *DataPlugin) Put(ctx core.RequestContext, id string, item Storable) error {
	return svc.DataComponent.Put(ctx, id, item)
}

//upsert an object ...insert if not there... update if there
func (svc *DataPlugin) UpsertId(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return svc.DataComponent.UpsertId(ctx, id, newVals)
}

func (svc *DataPlugin) Update(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	return svc.DataComponent.Update(ctx, id, newVals)
}

func (svc *DataPlugin) Upsert(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return svc.DataComponent.Upsert(ctx, queryCond, newVals)
}

func (svc *DataPlugin) UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	return svc.DataComponent.UpdateAll(ctx, queryCond, newVals)
}

//update objects by ids, fields to be updated should be provided as key value pairs
func (svc *DataPlugin) UpdateMulti(ctx core.RequestContext, ids []string, newVals map[string]interface{}) error {
	return svc.DataComponent.UpdateMulti(ctx, ids, newVals)
}

//item must support Deleted field for soft deletes
func (svc *DataPlugin) Delete(ctx core.RequestContext, id string) error {
	return svc.DataComponent.Delete(ctx, id)
}

//Delete object by ids
func (svc *DataPlugin) DeleteMulti(ctx core.RequestContext, ids []string) error {
	return svc.DataComponent.DeleteMulti(ctx, ids)
}

//Delete object by condition
func (svc *DataPlugin) DeleteAll(ctx core.RequestContext, queryCond interface{}) ([]string, error) {
	return svc.DataComponent.DeleteAll(ctx, queryCond)
}

func (svc *DataPlugin) GetById(ctx core.RequestContext, id string) (Storable, error) {
	return svc.DataComponent.GetById(ctx, id)
}

//Get multiple objects by id
func (svc *DataPlugin) GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]Storable, error) {
	return svc.DataComponent.GetMulti(ctx, ids, orderBy)
}

func (svc *DataPlugin) GetMultiHash(ctx core.RequestContext, ids []string) (map[string]Storable, error) {
	return svc.DataComponent.GetMultiHash(ctx, ids)
}

func (svc *DataPlugin) Count(ctx core.RequestContext, queryCond interface{}) (count int, err error) {
	return svc.DataComponent.Count(ctx, queryCond)
}

func (svc *DataPlugin) CountGroups(ctx core.RequestContext, queryCond interface{}, group string) (res map[string]interface{}, err error) {
	return svc.DataComponent.CountGroups(ctx, queryCond, group)
}

func (svc *DataPlugin) GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.DataComponent.GetList(ctx, pageSize, pageNum, mode, orderBy)
}

func (svc *DataPlugin) Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.DataComponent.Get(ctx, queryCond, pageSize, pageNum, mode, orderBy)
}

//create condition for passing to data service
func (svc *DataPlugin) CreateCondition(ctx core.RequestContext, operation ConditionType, args ...interface{}) (interface{}, error) {
	return svc.DataComponent.CreateCondition(ctx, operation, args...)
}
