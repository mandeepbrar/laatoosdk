package core

import (
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

type MetaDataProvider interface {
	CreateServiceInfo(name, description, version string, reqInfo RequestInfo, resInfo ResponseInfo, configurations []Configuration) ServiceInfo
	CreateFactoryInfo(name, description, version string, configurations []Configuration) ServiceFactoryInfo
	CreateModuleInfo(name, description, version string, configurations []Configuration) ModuleInfo
	CreateRequestInfo(params map[string]Param) RequestInfo
	CreateResponseInfo(params map[string]Param) ResponseInfo
	CreateConfiguration(name string, desc string, conftype datatypes.DataType, required bool, defaultValue interface{}, varToSet string) Configuration
	CreateParam(ctx ServerContext, name string, desc string, paramtype datatypes.DataType, customObjectType string, collectio, isStream bool, required bool) (Param, error)
}

type Info interface {
	GetDescription() string
	GetType() string
	GetVersion() string
	GetProperty(string) interface{}
}

type ConfigurableObjectInfo interface {
	Info
	GetConfigurations() map[string]Configuration
}

type ServiceInfo interface {
	ConfigurableObjectInfo
	GetRequestInfo() RequestInfo
	GetResponseInfo() ResponseInfo
	IsComponent() bool
}

type RequestInfo interface {
	ParamInfo() map[string]Param
}

type ResponseInfo interface {
	ParamInfo() map[string]Param
}

type ModuleInfo interface {
	ConfigurableObjectInfo
}

type ServiceFactoryInfo interface {
	ConfigurableObjectInfo
}

type defaultInfo struct {
	description string
	objtype     string
	objversion  string
	properties  utils.StringMap
}

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
