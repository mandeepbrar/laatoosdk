package core

import (
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

// MetaDataProvider is handed to a plugin's Manifest function so the plugin can build metadata
// objects of the server's internal types (laatooserver core/pluginloader.go:65-79).
//
// The values it returns are opaque and several of them are type-asserted back to those internal
// types, so a hand-rolled implementation of ServiceInfo, RequestInfo or ResponseInfo cannot be
// handed back in. No plugin under laatoomodules currently uses this interface; plugins declare
// metadata in registry/objects YAML instead.
type MetaDataProvider interface {
	// CreateServiceInfo builds a ServiceInfo. Its reqInfo handling is DEFECTIVE: the
	// implementation branches on resInfo when deciding what to do with reqInfo
	// (laatooserver core/serviceinfo.go:36-40), so passing a reqInfo together with a nil resInfo
	// silently discards the reqInfo and substitutes an empty request, and passing a nil reqInfo
	// together with a non-nil resInfo panics on a nil type assertion. Only the (nil, nil) case
	// behaves as written, and it is the only case the server's own test covers. reqInfo and
	// resInfo must be values returned by CreateRequestInfo / CreateResponseInfo.
	CreateServiceInfo(name, description, version string, reqInfo RequestInfo, resInfo ResponseInfo, configurations []Configuration) ServiceInfo
	// CreateFactoryInfo builds a ServiceFactoryInfo whose Info type string is "Factory".
	CreateFactoryInfo(name, description, version string, configurations []Configuration) ServiceFactoryInfo
	// CreateModuleInfo builds a ModuleInfo whose Info type string is "Module".
	CreateModuleInfo(name, description, version string, configurations []Configuration) ModuleInfo
	// CreateRequestInfo wraps a param map as a RequestInfo. The map is stored by reference, not
	// copied, and a nil map is stored as-is -- which makes ParamInfo return a nil map that panics
	// if anything later writes to it.
	CreateRequestInfo(params map[string]Param) RequestInfo
	// CreateResponseInfo wraps a param map as a ResponseInfo, stored by reference like
	// CreateRequestInfo.
	CreateResponseInfo(params map[string]Param) ResponseInfo
	// CreateConfiguration builds one Configuration declaration. varToSet names the struct field
	// the parsed value is injected into -- the role `variable:` plays in an object spec's
	// configurations block; pass "" for no injection. Unlike the Add*Configuration helpers on
	// ConfigurableObject, required is honoured exactly as given.
	CreateConfiguration(name string, desc string, conftype datatypes.DataType, required bool, defaultValue interface{}, varToSet string) Configuration
	// CreateParam builds one Param declaration. The error is ALWAYS nil
	// (laatooserver core/param.go:34-37); the call panics instead when no ObjectLoader server
	// element is registered, because it type-asserts that element without a comma-ok. Note the
	// argument order -- (collection, isStream, required) -- differs from Service.AddParam's
	// (collection, required, stream).
	CreateParam(ctx ServerContext, name string, desc string, paramtype datatypes.DataType, customObjectType string, collectio, isStream bool, required bool) (Param, error)
}

// Info is the common metadata carried by every configurable element and every registered object.
type Info interface {
	// GetDescription returns the description, taken from the object spec's `description` key or
	// from the constructor argument.
	GetDescription() string
	// GetType returns the object KIND, not a data type. The programmatic constructors set it to
	// "Service", "Factory" or "Module", while the YAML-driven path copies the spec's `type:` key
	// verbatim -- so a spec written `type: service` yields the lowercase string. Compare
	// case-insensitively.
	GetType() string
	// GetVersion returns the version string, from the object spec's `version` key or the
	// constructor argument. For entity objects registered from a generated
	// autogen_objectsmanifest.go it is frequently the literal "%!s(<nil>)", an artefact of the
	// code generator rather than a version.
	GetVersion() string
	// GetProperty returns a named extra property, or nil when it is absent or the Info was built
	// with no properties. Every Info the server constructs for a service, factory or module
	// carries an EMPTY property map and there is no setter, so GetProperty always returns nil for
	// those. It is populated only for objects registered from a plugin manifest, where the
	// generated manifest stores the entity descriptor JSON under the key "descriptor".
	GetProperty(string) interface{}
}

// ConfigurableObjectInfo is the metadata of an element that has declared configurations.
type ConfigurableObjectInfo interface {
	Info
	// GetConfigurations returns the live declaration map keyed by configuration name. Each
	// Configuration's GetValue is nil until the element's configuration has been processed.
	GetConfigurations() map[string]Configuration
}

// ServiceInfo is a service's metadata: its configurations plus its request and response shapes.
type ServiceInfo interface {
	ConfigurableObjectInfo
	// GetRequestInfo returns the declared request parameters. Never nil for a service the server
	// loaded, but its ParamInfo map may be empty when the object spec declares no request params.
	GetRequestInfo() RequestInfo
	// GetResponseInfo returns the declared response parameters. Nothing in laatooserver validates
	// a response against them; they exist for descriptor and tool generation.
	GetResponseInfo() ResponseInfo
	// IsComponent reports the flag set by Service.SetComponent or by the service YAML's
	// `component:` key. NOTHING READS IT: no code in laatooserver or laatoomodules calls this
	// method, so marking a service a component has no effect on registration, injection or
	// channel exposure.
	IsComponent() bool
}

// RequestInfo is the declared request parameter set of a service.
type RequestInfo interface {
	// ParamInfo returns the live declaration map, not a copy. The Params in it are declarations:
	// their GetValue is always nil (per-request values live on clones held by the Request).
	ParamInfo() map[string]Param
}

// ResponseInfo is the declared response parameter set of a service.
type ResponseInfo interface {
	// ParamInfo returns the live declaration map, not a copy.
	ParamInfo() map[string]Param
}

// ModuleInfo is a module's metadata. The server builds it from the module's object spec at
// src/server/registry/objects/<objectName>.yml, overwriting anything the object loader had
// already supplied (laatooserver core/module.go:127-131).
type ModuleInfo interface {
	ConfigurableObjectInfo
}

// ServiceFactoryInfo is a service factory's metadata.
type ServiceFactoryInfo interface {
	ConfigurableObjectInfo
}

type defaultInfo struct {
	description string
	objtype     string
	objversion  string
	properties  utils.StringMap
}

// NewInfo builds a plain Info. Mind the argument order: it is
// (description, objtype, objversion, props), NOT (name, objtype, description, ...). Callers get
// this wrong -- laatoomodules/ai/.../laatooai/src/server/go/manifest.go:13 passes a name as the
// description and a description as the version -- and nothing detects it, because all three are
// plain strings. props may be nil; GetProperty then returns nil for every key.
func NewInfo(description, objtype, objversion string, props utils.StringMap) Info {
	return &defaultInfo{description, objtype, objversion, props}
}

func (inf *defaultInfo) GetDescription() string {
	return inf.description
}
func (inf *defaultInfo) GetType() string {
	return inf.objtype
}
func (inf *defaultInfo) GetVersion() string {
	return inf.objversion
}

func (inf *defaultInfo) GetProperty(prop string) interface{} {
	if inf.properties != nil {
		return inf.properties[prop]
	}
	return nil
}

/*
func CreateServiceMetaData(description, requesttype string, params, configurations [][]string) interface{} {
	return map[string] interface{} { "Description": description, "RequestType": requesttype, "Params": params, "Configurations": configurations}
}


type Configuration struct {
	Name         string
	Conftype     string
	Required     string
	DefaultValue interface{}
}

type RequestInfo struct {
	DataType   string
	Collection string
	Stream     string
	Params     []Param
}

type ResponseInfo struct {
	Stream bool
}

type Param struct {
	Name       string
	Collection string
	DataType   string
}

type ServiceMetaData struct {
	Request        RequestInfo
	Response       ResponseInfo
	Description    string
	Component      string
	Configurations []Configuration
}

type ServiceFactoryMetaData struct {
	Description    string
	Configurations []Configuration
}

type ModuleMetaData struct {
	Description    string
	Configurations []Configuration
}



func CreateFactoryMetaData(description string, configurations [][]string) *ServiceFactoryMetaData {
	metadata := &ServiceFactoryMetaData{Description: description}
	configurationsCollection := make([]Configuration, len(configurations))
	for ind, confrow := range configurations {
		if len(confrow) < 4 {
			return nil
		}
		configurationsCollection[ind] = Configuration{confrow[0], confrow[1], confrow[2], confrow[3]}
	}
	metadata.Configurations = configurationsCollection
	return metadata
}

func CreateModuleMetaData(description string, configurations [][]string) *ModuleMetaData {
	metadata := &ModuleMetaData{Description: description}
	configurationsCollection := make([]Configuration, len(configurations))
	for ind, confrow := range configurations {
		if len(confrow) < 4 {
			return nil
		}
		configurationsCollection[ind] = Configuration{confrow[0], confrow[1], confrow[2], confrow[3]}
	}
	metadata.Configurations = configurationsCollection
	return metadata
}
*/
