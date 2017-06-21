package common

/*
import (
	"laatoo/sdk/components/data"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
	"laatoo/sdk/log"
)

type deleteRefOperation struct {
	*refOperation
	do func(ctx core.RequestContext, ids []string) error
}

func buildCascadedDeleteOperation(ctx core.ServerContext, conf config.Config, opname string, targetsvcname string) (RefOperation, error) {
	targetfield, ok := conf.GetString(CONF_REF_TARG_FIELD)
	if !ok {
		return nil, errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Conf", CONF_REF_TARG_FIELD, "Operation", opname)
	}
	opr := &deleteRefOperation{refOperation: &refOperation{name: opname, targetsvcname: targetsvcname}}
	opr.do = func(newctx core.RequestContext, ids []string) error {
		return cascadeDelete(newctx, opr.targetService, targetfield, ids)
	}
	return opr, nil
}

func cascadeDelete(ctx core.RequestContext, dataService data.DataComponent, targetfield string, ids []string) error {
	if dataService.Supports(data.InQueries) {
		condition, _ := dataService.CreateCondition(ctx, data.MATCHMULTIPLEVALUES, targetfield, ids)
		_, err := dataService.DeleteAll(ctx, condition)
		return err
	} else {
		for _, id := range ids {
			condition, _ := dataService.CreateCondition(ctx, data.FIELDVALUE, map[string]interface{}{targetfield: id})
			_, err := dataService.DeleteAll(ctx, condition)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteRefOps(ctx core.RequestContext, opers []RefOperation, ids []string) error {
	if opers != nil {
		log.Trace(ctx, "deleterefops")
		for _, oper := range opers {
			dr := oper.(*deleteRefOperation)
			log.Trace(ctx, "deleterefops", "oper", dr.name)
			err := dr.do(ctx, ids)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
*/
