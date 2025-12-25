package errors

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/log"
)

const (
	CORE_ERROR_WRAPPER             = "Wrapper"
	CORE_ERROR_PROVIDER_NOT_FOUND  = "Core_Provider_Not_Found"
	CORE_ERROR_CODEC_NOT_FOUND     = "Core_Codec_Not_Found"
	CORE_ERROR_MISSING_SERVICE     = "Core_Missing_Service"
	CORE_ERROR_BAD_ARG             = "Core_Bad_Arg"
	CORE_ERROR_BAD_REQUEST         = "Core_Bad_Request"
	CORE_ERROR_MISSING_ARG         = "Core_Missing_Arg"
	CORE_ERROR_MISSING_CONF        = "Core_Missing_Conf"
	CORE_ERROR_MISSING_PLUGIN      = "Core_Missing_Plugin"
	CORE_ERROR_BAD_CONF            = "Core_Bad_Conf"
	CORE_ERROR_UNAUTHORIZED        = "Core_Error_Unauthorized"
	CORE_ERROR_RES_NOT_FOUND       = "Core_Resource_Not_Found"
	CORE_ERROR_DEP_NOT_MET         = "Core_Dep_Not_Met"
	CORE_ERROR_TYPE_MISMATCH       = "Core_Type_Mismatch"
	CORE_ERROR_NOT_IMPLEMENTED     = "Core_Not_Implemented"
	CORE_ERROR_PLUGIN_NOT_LOADED   = "Core_Plugin_Not_Loaded"
	CORE_ERROR_TENANT_MISMATCH     = "Core_Tenant_Mismatch"
	CORE_ERROR_INVALID_PAYLOAD     = "Core_Invalid_Payload"
	CORE_ERROR_INTERNAL_ERROR      = "Core_Internal_Error"
	CORE_ERROR_SERIALIZATION_ERROR = "Core_Serialization_Error"
)

func init() {
	RegisterCode(CORE_ERROR_WRAPPER, "Wrapped error.")
	RegisterCode(CORE_ERROR_PROVIDER_NOT_FOUND, "Factory not registered.")
	RegisterCode(CORE_ERROR_CODEC_NOT_FOUND, "Codec not registered.")
	RegisterCode(CORE_ERROR_MISSING_SERVICE, "Expected service is missing.")
	RegisterCode(CORE_ERROR_PLUGIN_NOT_LOADED, "Plugins could not be loaded.")
	RegisterCode(CORE_ERROR_MISSING_ARG, "All arguments have not been provided for the call.")
	RegisterCode(CORE_ERROR_MISSING_PLUGIN, "Required plugin is missing.")
	RegisterCode(CORE_ERROR_BAD_ARG, "Invalid argument was provided.")
	RegisterCode(CORE_ERROR_BAD_REQUEST, "Invalid request was sent.")
	RegisterCode(CORE_ERROR_MISSING_CONF, "All configurations have not been provided.")
	RegisterCode(CORE_ERROR_BAD_CONF, "Configuration is not correct.")
	RegisterCode(CORE_ERROR_UNAUTHORIZED, "You are not allowed to access this resource.")
	RegisterCode(CORE_ERROR_RES_NOT_FOUND, "Requested resource was not found.")
	RegisterCode(CORE_ERROR_DEP_NOT_MET, "Dependency could not be met.")
	RegisterCode(CORE_ERROR_TYPE_MISMATCH, "Type Mismatch.")
	RegisterCode(CORE_ERROR_NOT_IMPLEMENTED, "Method has not been implemented by this service.")
	RegisterCode(CORE_ERROR_TENANT_MISMATCH, "Tenant Mismatch.")
	RegisterCode(CORE_ERROR_INTERNAL_ERROR, "Internal Error")
	RegisterCode(CORE_ERROR_SERIALIZATION_ERROR, "Serialization error")
}

func WrapError(ctx ctx.Context, err error, info ...interface{}) error {
	if err != nil {
		laatooErr, ok := err.(*Error)
		if ok {
			log.Debug(ctx, laatooErr.error.Error(), append(laatooErr.info, info...)...)
			return err
		} else {
			return RethrowError(ctx, CORE_ERROR_WRAPPER, err, info...)
		}
	}
	return nil
}

func WrapErrorWithCode(ctx ctx.Context, err error, errCode string, info ...interface{}) error {
	if err != nil {
		_, ok := err.(*Error)
		if ok {
			return err
		} else {
			return RethrowError(ctx, errCode, err, info...)
		}
	}
	return nil
}

func BadRequest(ctx ctx.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_REQUEST, info...)
}

func BadArg(ctx ctx.Context, argName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_ARG, append(info, "Argument", argName)...)
}

func MissingArg(ctx ctx.Context, argName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_ARG, append(info, "Argument", argName)...)
}

func BadConf(ctx ctx.Context, confName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_BAD_CONF, append(info, "Configuration", confName)...)
}

func DepNotMet(ctx ctx.Context, dep string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_DEP_NOT_MET, append(info, "Dependency", dep)...)
}

func MissingConf(ctx ctx.Context, confName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_CONF, append(info, "Configuration", confName)...)
}

func MissingService(ctx ctx.Context, svcName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_MISSING_SERVICE, append(info, "Service", svcName)...)
}

func NotImplemented(ctx ctx.Context, methodName string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_NOT_IMPLEMENTED, append(info, "Method", methodName)...)
}

func NotFound(ctx ctx.Context, resource string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_RES_NOT_FOUND, append(info, "Resource", resource)...)
}
func TypeMismatch(ctx ctx.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_TYPE_MISMATCH, info...)
}

func Unauthorized(ctx ctx.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_UNAUTHORIZED, info...)
}

func InternalError(ctx ctx.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_INTERNAL_ERROR, info...)
}

func InvalidPayload(ctx ctx.Context, key string, errorReason string, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_INVALID_PAYLOAD, append(info, "Key", key, "Error Reason", errorReason)...)
}
func SerializationError(ctx ctx.Context, info ...interface{}) error {
	return ThrowError(ctx, CORE_ERROR_SERIALIZATION_ERROR, info...)
}
