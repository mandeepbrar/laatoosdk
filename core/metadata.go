package core

type MetaDataProvider interface {
	CreateServiceInfo(description string, reqInfo RequestInfo, streamedResponse bool, configurations []Configuration) ServiceInfo
	CreateFactoryInfo(description string, configurations []Configuration) ServiceFactoryInfo
	CreateModuleInfo(description string, configurations []Configuration) ModuleInfo
	CreateRequestInfo(requesttype string, collection bool, stream bool, params []Param) RequestInfo
	CreateConfiguration(name, conftype string, required bool, defaultValue interface{}) Configuration
	CreateParam(name, paramtype string, collectio bool, required bool) Param
}

type Info interface {
	GetDescription() string
	GetType() string
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
	GetRequiredServices() map[string]string
}

type RequestInfo interface {
	GetDataType() string
	IsCollection() bool
	IsStream() bool
	GetParams() map[string]Param
}

type ResponseInfo interface {
	IsStream() bool
}

type ModuleInfo interface {
	ConfigurableObjectInfo
}

type ServiceFactoryInfo interface {
	ConfigurableObjectInfo
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
