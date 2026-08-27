package core

import (
	"laatoo.io/sdk/config"
)

// ServiceFactoryProvider builds a ServiceFactory from configuration.
type ServiceFactoryProvider func(ctx ServerContext, config config.Config) (ServiceFactory, error)

// ServiceFactory creates service instances. A plugin factory struct embeds core.ServiceFactory
// and the server reflects a field literally named "ServiceFactory" and assigns its own
// implementation into it (laatooserver core/factory.go:39-48); a struct that does not embed
// core.ServiceFactory under that exact field name fails to load with Core_Error_Type_Mismatch.
//
// Lifecycle, in the order the server drives it:
//
//	Describe -> factory config processed against what Describe declared
//	  -> configuration values injected into struct fields -> Initialize -> Start
//	  -> CreateService, once per service the factory owns -> Stop -> Unload
type ServiceFactory interface {
	ConfigurableObject
	// Metadata returns the factory's ServiceFactoryInfo -- its declared configurations,
	// description and version, seeded from the factory's object spec.
	Metadata() ServiceFactoryInfo
	// Describe is where a factory declares its configurations, before any configuration is
	// parsed. Its error is DISCARDED -- the server calls it as a bare statement
	// (laatooserver core/factory.go:60) -- so returning an error here neither aborts the load nor
	// gets logged. This differs from Service.Describe, whose error IS checked. Signal failure
	// from Initialize instead.
	Describe(ServerContext) error
	// Initialize runs after configuration has been parsed against the declarations Describe made
	// and after configuration-backed struct fields have been injected. conf is the raw factory
	// config.
	Initialize(ctx ServerContext, conf config.Config) error
	// Start runs after Initialize, before any CreateService call.
	Start(ctx ServerContext) error
	// Stop is called during factory unload, before Unload. An error returned here aborts the
	// unload and Unload is not called (laatooserver core/factorymanager.go:293-300).
	Stop(ctx ServerContext) error
	// Unload is called immediately after a successful Stop.
	Unload(ctx ServerContext) error
	// CreateService returns a new service instance. name is the service alias and method is the
	// registered object name from the service YAML's `servicemethod:`; conf is that service's
	// configuration.
	//
	// The base implementation the server injects returns (nil, nil)
	// (laatooserver core/factoryimpl.go:61-63), so a factory that forgets to override this
	// produces no error of its own. The caller does guard on the nil service and raises
	// Core_Error_Missing_Service naming the alias (core/servicemanager.go:461-467) -- so the
	// symptom is a missing-service failure at startup rather than a nil dereference, and the
	// error names the service, not the factory that failed to build it.
	CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error)
	// ServerElement returns the ServerElementServiceFactory proxy wrapping this factory, for code
	// that needs the factory's name, context or owning module.
	ServerElement() ServerElement
}
