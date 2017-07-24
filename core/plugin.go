package core

type PluginComponent struct {
	Name                    string
	Object                  interface{}
	ObjectCollectionCreator ObjectCollectionCreator
	ObjectCreator           ObjectCreator
	ObjectFactory           ObjectFactory
}

//manifest that needs to be provided by every plugin
type PluginManifest func() []PluginComponent
