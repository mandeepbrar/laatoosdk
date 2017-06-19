package common

/*
import (
	"laatoo/sdk/components/data"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
	"laatoo/sdk/log"
)

type getRefOperation struct {
	*refOperation
	do func(ctx core.RequestContext, ids []string, inputData interface{}) (interface{}, error)
}

func buildJoinOperation(ctx core.ServerContext, conf config.Config, opname string, targetsvcname string) (RefOperation, error) {
	targetfield, ok := conf.GetString(CONF_REF_TARG_FIELD)
	if !ok {
		return nil, errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Conf", CONF_REF_TARG_FIELD, "Operation", opname)
	}
	opr := &getRefOperation{refOperation: &refOperation{name: opname, targetsvcname: targetsvcname}}
	opr.do = func(newctx core.RequestContext, ids []string, inputdata interface{}) (interface{}, error) {
		return fillJoin(newctx, opr.targetService, targetfield, ids, inputdata)
	}
	return opr, nil
}

func fillJoin(ctx core.RequestContext, dataService data.DataComponent, targetfield string, ids []string, inputData interface{}) (interface{}, error) {
	ctx = ctx.SubContext("Join")
	hash, err := dataService.GetMultiHash(ctx, ids)
	if err != nil {
		return nil, err
	}
	arrayData, ok := inputData.([]data.Storable)
	if ok {
		for _, stor := range arrayData {
			id := stor.GetId()
			joinedItem, ok := hash[id]
			if ok {
				log.Logger.Info(ctx, "Joining item", "item", joinedItem)
				stor.Join(joinedItem)
			}
		}
	} else {
		mapData, _ := inputData.(map[string]data.Storable)
		for id, stor := range mapData {
			joinedItem, ok := hash[id]
			if ok {
				stor.Join(joinedItem)
			}
		}
	}
	return inputData, nil
}

func GetRefOps(ctx core.RequestContext, opers []RefOperation, ids []string, inputData interface{}) (interface{}, error) {
	if opers != nil {
		var err error
		log.Logger.Trace(ctx, "getrefops")
		for _, oper := range opers {
			gr := oper.(*getRefOperation)
			log.Logger.Trace(ctx, "getrefops", "oper", gr.name)
			inputData, err = gr.do(ctx, ids, inputData)
			if err != nil {
				return nil, err
			}
		}
	}
	return inputData, nil
}
*/
