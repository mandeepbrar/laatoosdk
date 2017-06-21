package common

/*
import (
	"laatoo/sdk/components/data"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
	"laatoo/sdk/log"
	"laatoo/sdk/utils"
	"reflect"
)

type cachedDataService struct {
	*data.BaseComponent
	rootComp                data.DataComponent
	rootService             core.Service
	object                  string
	objectCollectionCreator core.ObjectCollectionCreator
	objectCreator           core.ObjectCreator
	objType                 reflect.Type
}

func NewCachedDataService(ctx core.ServerContext, svc data.DataComponent) *cachedDataService {
	return &cachedDataService{rootComp: svc, rootService: svc.(core.Service)}
}

func (svc *cachedDataService) Initialize(ctx core.ServerContext, conf config.Config) error {
	object, ok := conf.GetString(CONF_DATA_OBJECT)
	if !ok {
		return errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Missing Conf", CONF_DATA_OBJECT)
	}
	objectCreator, err := ctx.GetObjectCreator(object)
	if err != nil {
		return errors.ThrowError(ctx, errors.CORE_ERROR_BAD_ARG, "Could not get Object creator for", object)
	}
	objectCollectionCreator, err := ctx.GetObjectCollectionCreator(object)
	if err != nil {
		return errors.ThrowError(ctx, errors.CORE_ERROR_BAD_ARG, "Could not get Object Collection creator for", object)
	}

	testObj := objectCreator()
	svc.objType = reflect.TypeOf(testObj)

	svc.object = object
	svc.objectCreator = objectCreator
	svc.objectCollectionCreator = objectCollectionCreator
	return svc.rootService.Initialize(ctx, conf)
}

func (svc *cachedDataService) Start(ctx core.ServerContext) error {
	return svc.rootService.Start(ctx)
}

func (svc *cachedDataService) Invoke(ctx core.RequestContext) error {
	return nil
}

/*
func (svc *cachedDataService) CreateDBCollection(ctx core.RequestContext) error {
	return svc.rootComp.CreateDBCollection(ctx)
}

func (svc *cachedDataService) DropCollection(ctx core.RequestContext) error {
	return svc.rootComp.DropCollection(ctx)
}

func (svc *cachedDataService) CollectionExists(ctx core.RequestContext) (bool, error) {
	return svc.rootComp.CollectionExists(ctx)
}

func (svc *cachedDataService) GetDataServiceType() string {
	return svc.rootComp.GetDataServiceType()
}

func (svc *cachedDataService) GetObject() string {
	return svc.rootComp.GetObject()
}

func (svc *cachedDataService) Supports(feature data.Feature) bool {
	return svc.rootComp.Supports(feature)
}

func (svc *cachedDataService) Save(ctx core.RequestContext, item data.Storable) error {
	return svc.rootComp.Save(ctx, item)
}

func (svc *cachedDataService) PutMulti(ctx core.RequestContext, items []data.Storable) error {
	err := svc.rootComp.PutMulti(ctx, items)
	if err != nil {
		for _, item := range items {
			id := item.GetId()
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return err
}

func (svc *cachedDataService) Put(ctx core.RequestContext, id string, item data.Storable) error {
	err := svc.rootComp.Put(ctx, id, item)
	if err != nil {
		ctx.InvalidateCache(svc.object, id)
	}
	return err
}

//upsert an object ...insert if not there... update if there
func (svc *cachedDataService) UpsertId(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	err := svc.rootComp.UpsertId(ctx, id, newVals)
	if err != nil {
		ctx.InvalidateCache(svc.object, id)
	}
	return err
}

func (svc *cachedDataService) Update(ctx core.RequestContext, id string, newVals map[string]interface{}) error {
	err := svc.rootComp.Update(ctx, id, newVals)
	if err != nil {
		ctx.InvalidateCache(svc.object, id)
	}
	return err
}

func (svc *cachedDataService) Upsert(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	ids, err := svc.rootComp.Upsert(ctx, queryCond, newVals)
	if err != nil {
		for _, id := range ids {
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return ids, err
}

func (svc *cachedDataService) UpdateAll(ctx core.RequestContext, queryCond interface{}, newVals map[string]interface{}) ([]string, error) {
	ids, err := svc.rootComp.UpdateAll(ctx, queryCond, newVals)
	if err != nil {
		for _, id := range ids {
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return ids, err
}

//update objects by ids, fields to be updated should be provided as key value pairs
func (svc *cachedDataService) UpdateMulti(ctx core.RequestContext, ids []string, newVals map[string]interface{}) error {
	err := svc.rootComp.UpdateMulti(ctx, ids, newVals)
	if err != nil {
		for _, id := range ids {
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return err
}

//item must support Deleted field for soft deletes
func (svc *cachedDataService) Delete(ctx core.RequestContext, id string) error {
	err := svc.rootComp.Delete(ctx, id)
	if err != nil {
		ctx.InvalidateCache(svc.object, id)
	}
	return err
}

//Delete object by ids
func (svc *cachedDataService) DeleteMulti(ctx core.RequestContext, ids []string) error {
	err := svc.rootComp.DeleteMulti(ctx, ids)
	if err != nil {
		for _, id := range ids {
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return err
}

//Delete object by condition
func (svc *cachedDataService) DeleteAll(ctx core.RequestContext, queryCond interface{}) ([]string, error) {
	ids, err := svc.rootComp.DeleteAll(ctx, queryCond)
	if err != nil {
		for _, id := range ids {
			ctx.InvalidateCache(svc.object, id)
		}
	}
	return ids, err
}

func (svc *cachedDataService) GetById(ctx core.RequestContext, id string) (data.Storable, error) {
	ctx = ctx.SubContext("Cache_GetById")

	ent, ok := ctx.GetObjectFromCache(svc.object, id, svc.object)
	if ok {
		log.Trace(ctx, "Cache hit")
		return ent.(data.Storable), nil
	}
	stor, err := svc.rootComp.GetById(ctx, id)
	if err == nil {
		ctx.PutInCache(svc.object, id, stor)
	}
	return stor, err
}

//Get multiple objects by id
func (svc *cachedDataService) GetMulti(ctx core.RequestContext, ids []string, orderBy string) ([]data.Storable, error) {
	ctx = ctx.SubContext("Cache_GetMulti")
	res := make([]data.Storable, len(ids))
	cachedItems := ctx.GetObjectsFromCache(svc.object, ids, svc.object)
	idsNotCached := make([]string, 0, 10)
	for index, id := range ids {
		item, ok := cachedItems[id]
		if !ok || item == nil {
			idsNotCached = append(idsNotCached, id)
		} else {
			res[index] = item.(data.Storable)
		}
	}
	if len(idsNotCached) == 0 {
		return res, nil
	}
	stormap, err := svc.rootComp.GetMultiHash(ctx, idsNotCached)
	if err == nil {
		for index, id := range ids {
			if res[index] == nil {
				item, ok := stormap[id]
				if ok {
					ctx.PutInCache(svc.object, id, item)
					res[index] = item
				}
			}
		}
	}
	log.Error(ctx, "Elapsed time after db", "time", ctx.GetElapsedTime())
	return res, err
}

func (svc *cachedDataService) GetMultiHash(ctx core.RequestContext, ids []string, orderBy string) (map[string]data.Storable, error) {
	ctx = ctx.SubContext("Cache_GetMultiHash")
	res, err := svc.rootComp.GetMultiHash(ctx, ids)
	if err == nil {
		ctx.PutMultiInCache(svc.object, utils.CastToStringMap(res))
	}
	return res, err
}

func (svc *cachedDataService) Count(ctx core.RequestContext, queryCond interface{}) (count int, err error) {
	return svc.rootComp.Count(ctx, queryCond)
}

func (svc *cachedDataService) CountGroups(ctx core.RequestContext, queryCond interface{}, group string) (res map[string]interface{}, err error) {
	return svc.rootComp.CountGroups(ctx, queryCond, group)
}

func (svc *cachedDataService) GetList(ctx core.RequestContext, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []data.Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.rootComp.GetList(ctx, pageSize, pageNum, mode, orderBy)
}

func (svc *cachedDataService) Get(ctx core.RequestContext, queryCond interface{}, pageSize int, pageNum int, mode string, orderBy string) (dataToReturn []data.Storable, ids []string, totalrecs int, recsreturned int, err error) {
	return svc.rootComp.Get(ctx, queryCond, pageSize, pageNum, mode, orderBy)
}

//create condition for passing to data service
func (svc *cachedDataService) CreateCondition(ctx core.RequestContext, operation data.ConditionType, args ...interface{}) (interface{}, error) {
	return svc.rootComp.CreateCondition(ctx, operation, args...)
}*/
