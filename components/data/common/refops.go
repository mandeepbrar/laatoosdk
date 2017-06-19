package common

/*
import (
	"laatoo/sdk/components/data"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/errors"
)

const (
	CONF_REF_OPS        = "reference_operations"
	CONF_REF_OP         = "operation"
	CONF_REF_TARG_SVC   = "targetsvc"
	CONF_REF_TARG_FIELD = "targetfield"
)

type RefOperation interface {
	Initialize(ctx core.ServerContext) error
}

type refOperation struct {
	name          string
	targetsvcname string
	targetService data.DataComponent
}

func BuildRefOps(ctx core.ServerContext, conf config.Config) (deleterefops []RefOperation, updaterefops []RefOperation, saverefops []RefOperation, getrefops []RefOperation, err error) {
	deleteOps := []RefOperation{}
	updaterefOps := []RefOperation{}
	saverefOps := []RefOperation{}
	getrefOps := []RefOperation{}
	refOps, ok := conf.GetSubConfig(CONF_REF_OPS)
	if ok {
		opnames := refOps.AllConfigurations()
		for _, opname := range opnames {
			operConf, _ := refOps.GetSubConfig(opname)
			operation, ok := operConf.GetString(CONF_REF_OP)
			if !ok {
				return nil, nil, nil, nil, errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Conf", CONF_REF_OP, "Operation", opname)
			}
			targetsvc, ok := operConf.GetString(CONF_REF_TARG_SVC)
			if !ok {
				return nil, nil, nil, nil, errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Conf", CONF_REF_TARG_SVC, "Operation", opname)
			}
			switch operation {
			case "CascadedDelete":
				oper, err := buildCascadedDeleteOperation(ctx, operConf, opname, targetsvc)
				if err != nil {
					return nil, nil, nil, nil, errors.WrapError(ctx, err)
				}
				deleteOps = append(deleteOps, oper)
				break
			case "Join":
				oper, err := buildJoinOperation(ctx, operConf, opname, targetsvc)
				if err != nil {
					return nil, nil, nil, nil, errors.WrapError(ctx, err)
				}
				getrefOps = append(getrefOps, oper)
				break
			case "Save":
			case "Update":

			}
		}
	}
	return deleteOps, updaterefOps, saverefOps, getrefOps, nil
}

func InitialRefOps(ctx core.ServerContext, arr []RefOperation) error {
	if arr != nil {
		for _, refop := range arr {
			err := refop.Initialize(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (oper *refOperation) Initialize(ctx core.ServerContext) error {
	svc, err := ctx.GetService(oper.targetsvcname)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	targetsvc, ok := svc.(data.DataComponent)
	if !ok {
		return errors.ThrowError(ctx, errors.CORE_ERROR_BAD_CONF)
	}
	oper.targetService = targetsvc
	return nil
}*/
