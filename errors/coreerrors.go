package errors

import (
	"fmt"
	"laatoo/sdk/core"
)

const (
	CORE_ERROR_WRAPPER            = "Wrapper"
	CORE_ERROR_PROVIDER_NOT_FOUND = "Core_Provider_Not_Found"
	CORE_ERROR_CODEC_NOT_FOUND    = "Core_Codec_Not_Found"
	CORE_ERROR_MISSING_SERVICE    = "Core_Missing_Service"
	CORE_ERROR_BAD_ARG            = "Core_Bad_Arg"
	CORE_ERROR_BAD_REQUEST        = "Core_Bad_Request"
	CORE_ERROR_MISSING_ARG        = "Core_Missing_Arg"
	CORE_ERROR_MISSING_CONF       = "Core_Missing_Conf"
	CORE_ERROR_MISSING_MODULE     = "Core_Missing_Module"
	CORE_ERROR_BAD_CONF           = "Core_Bad_Conf"
	CORE_ERROR_UNAUTHORIZED       = "Core_Error_Unauthorized"
	CORE_ERROR_RES_NOT_FOUND      = "Core_Resource_Not_Found"
	CORE_ERROR_DEP_NOT_MET        = "Core_Dep_Not_Met"
	CORE_ERROR_TYPE_MISMATCH      = "Core_Type_Mismatch"
	CORE_ERROR_NOT_IMPLEMENTED    = "Core_Not_Implemented"
	CORE_ERROR_PLUGIN_NOT_LOADED  = "Core_Plugin_Not_Loaded"
)

func init() {
	RegisterCode(CORE_ERROR_WRAPPER, FATAL, fmt.Errorf("Wrapped error."))
	RegisterCode(CORE_ERROR_PROVIDER_NOT_FOUND, FATAL, fmt.Errorf("Factory not registered."))
	RegisterCode(CORE_ERROR_CODEC_NOT_FOUND, FATAL, fmt.Errorf("Codec not registered."))
	RegisterCode(CORE_ERROR_MISSING_SERVICE, FATAL, fmt.Errorf("Expected service is missing."))
	RegisterCode(CORE_ERROR_PLUGIN_NOT_LOADED, FATAL, fmt.Errorf("Plugins could not be loaded."))
	RegisterCode(CORE_ERROR_MISSING_ARG, FATAL, fmt.Errorf("All arguments have not been provided for the call."))
	RegisterCode(CORE_ERROR_MISSING_MODULE, FATAL, fmt.Errorf("Required module is missing."))
	RegisterCode(CORE_ERROR_BAD_ARG, FATAL, fmt.Errorf("Invalid argument was provided."))
	RegisterCode(CORE_ERROR_BAD_REQUEST, FATAL, fmt.Errorf("Invalid request was sent."))
	RegisterCode(CORE_ERROR_MISSING_CONF, FATAL, fmt.Errorf("All configurations have not been provided."))
	RegisterCode(CORE_ERROR_BAD_CONF, FATAL, fmt.Errorf("Configuration is not correct."))
	RegisterCode(CORE_ERROR_UNAUTHORIZED, FATAL, fmt.Errorf("You are not allowed to access this resource."))
	RegisterCode(CORE_ERROR_RES_NOT_FOUND, FATAL, fmt.Errorf("Requested resource was not found."))
	RegisterCode(CORE_ERROR_DEP_NOT_MET, FATAL, fmt.Errorf("Dependency could not be met."))
	RegisterCode(CORE_ERROR_TYPE_MISMATCH, FATAL, fmt.Errorf("Type Mismatch."))
	RegisterCode(CORE_ERROR_NOT_IMPLEMENTED, FATAL, fmt.Errorf("Method has not been implemented by this service."))
}

func WrapError(ctx core.Context, err error, info ...interface{}) error {
	if err != nil {
		_, ok := err.(Error)
		if ok {
			return err
		} else {
			return RethrowError(ctx, CORE_ERROR_WRAPPER, err, info...)
		}
	}
	return nil
}

func BadRequest(ctx core.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_REQUEST, info...)
}

func BadArg(ctx core.Context, argName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_ARG, append(info, "Argument", argName)...)
}

func MissingArg(ctx core.Context, argName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_ARG, append(info, "Argument", argName)...)
}

func BadConf(ctx core.Context, confName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_CONF, append(info, "Configuration", confName)...)
}

func DepNotMet(ctx core.Context, dep string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_DEP_NOT_MET, append(info, "Dependency", dep)...)
}

func MissingConf(ctx core.Context, confName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_CONF, append(info, "Configuration", confName)...)
}

func MissingService(ctx core.Context, svcName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_SERVICE, append(info, "Service", svcName)...)
}

func NotImplemented(ctx core.Context, methodName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_NOT_IMPLEMENTED, append(info, "Method", methodName)...)
}

func TypeMismatch(ctx core.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_TYPE_MISMATCH, info...)
}

func Unauthorized(ctx core.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_UNAUTHORIZED, info...)
}
