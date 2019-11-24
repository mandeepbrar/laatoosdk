package core

type PluginComponent struct {
	Object        interface{}
	ObjectFactory ObjectFactory
	Metadata      Info
}

//manifest that needs to be provided by every module
type PluginManifest func(MetaDataProvider) []PluginComponent
