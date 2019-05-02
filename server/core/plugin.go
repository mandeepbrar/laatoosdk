package core

type PluginComponent struct {
	Name                    string
	Object                  interface{}
	ObjectCollectionCreator ObjectCollectionCreator
	ObjectCreator           ObjectCreator
	ObjectFactory           ObjectFactory
	Metadata                Info
}

//manifest that needs to be provided by every module
type PluginManifest func(MetaDataProvider) []PluginComponent
