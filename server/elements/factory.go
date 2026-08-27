package elements

import (
	"laatoo.io/sdk/server/core"
)

// Factory is the server's handle on a configured service factory -- the element the factory manager
// stores and the service manager asks to build each service instance.
//
// It wraps, but is not, the plugin's own core.ServiceFactory implementation -- Factory() reaches
// that.
// The context a factory was created in is significant beyond logging -- when a factory belongs to a
// lower level than the service being created, the service's context is overridden with the creating
// level's elements so the service sees the nearer solution/application/isolation.
type Factory interface {
	core.ServerElement

	// Factory returns the plugin's core.ServiceFactory implementation behind this element. This is
	// what the service manager calls CreateService on.
	Factory() core.ServiceFactory

	// GetModule returns the module that supplied this factory, or nil for a factory registered by
	// the server itself -- DefaultFactory and the activity factory both return nil.
	GetModule() core.Module

	// Metadata returns the factory's descriptor: its declared configurations, description and
	// version, loaded from the object spec registered for its Go type and then amended by the
	// factory's own Describe hook.
	//
	// The name on the returned info is overwritten with the factory's ALIAS, not the name in the
	// object spec, so two aliases of one factory type report distinct names.
	Metadata() core.ServiceFactoryInfo
}
