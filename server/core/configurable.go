package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
)

// Configuration is one declared configuration key on a service, factory or module -- what an
// object spec's `configurations:` block, or an Add*Configuration call inside Describe, produces.
type Configuration interface {
	// GetName returns the configuration key.
	GetName() string
	// GetDescription returns the declared description. It is always "" for a key declared through
	// AddStringConfiguration, which drops its desc argument -- see that method.
	GetDescription() string
	// IsRequired reports whether the key must be present in the element's configuration. A
	// missing required key fails startup with Core_Missing_Conf
	// (laatooserver core/configurableobject.go:289-292). Both AddStringConfiguration and
	// AddConfiguration set this true even when they are handed a default value.
	IsRequired() bool
	// GetDefaultValue returns the declared default. It is what GetConfiguration falls back to.
	// It is NOT applied during configuration parsing, so an absent required key still fails
	// before any caller can reach the default.
	GetDefaultValue() interface{}
	// GetValue returns the value parsed from configuration, or nil when the key was absent. It is
	// nil until the element's configuration has been processed, and therefore nil throughout
	// Describe.
	GetValue() interface{}
	// GetType returns the declared data type. It decides which config accessor parses the value
	// and hence the value's concrete Go type, which in turn decides which of the Get*Configuration
	// accessors below can read it back.
	GetType() datatypes.DataType
}

// ConfigurableObject is the configuration half of a Service, ServiceFactory or Module. The server
// injects its own implementation into the embedded core.Service / core.ServiceFactory /
// core.Module field, so these are the server's methods, not the plugin's.
//
// Declare in Describe; read from Initialize onwards. Two rules cause most of the surprises here.
//
//  1. A declared default does NOT make a key optional. AddStringConfiguration and
//     AddConfiguration mark the key required=true regardless of the default they are given
//     (laatooserver core/configurableobject.go:174-181), so an absent key fails startup with
//     Core_Missing_Conf and the default is never reached. Use AddOptionalConfiguration for
//     anything that has a usable default.
//
//  2. The bool returned by every Get*Configuration means "an explicit value was configured", NOT
//     "a value is available". When a key falls back to its default the value is returned with
//     ok == false (core/configurableobject.go:196-201), so the idiomatic
//     `if v, ok := ...; ok` silently discards every default. Read the value and ignore the flag,
//     or compare against the zero value.
//
// Reading with the wrong accessor is not uniformly safe: GetStringConfiguration and
// GetBoolConfiguration type-assert without a comma-ok and PANIC on a mismatch, while
// GetStringArrayConfiguration, GetStringsMapConfiguration and GetMapConfiguration return false.
type ConfigurableObject interface {
	// GetName returns the element's configured name -- the service alias, factory name or module
	// name. The server writes it over whatever the object spec carried.
	GetName() string
	// GetDescription returns the description from the object spec, overridden by any
	// Service.SetDescription made in Describe.
	GetDescription() string
	// GetVersion returns the object spec's `version` key, or "" when the spec omits it.
	GetVersion() string
	// GetConfigurations returns the live declaration map keyed by configuration name.
	GetConfigurations() map[string]Configuration
	//AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	// AddStringConfiguration declares a REQUIRED string configuration key, and DISCARDS desc: the
	// implementation passes "" as the description whatever is given
	// (laatooserver core/configurableobject.go:174-176). defaultValue is stored but unreachable in
	// the ordinary case, because required=true makes the element fail startup with
	// Core_Missing_Conf when the key is absent. For an optional string call
	// AddOptionalConfiguration(ctx, name, desc, datatypes.String, defaultValue) instead.
	AddStringConfiguration(ctx ServerContext, name string, desc string, defaultValue string)
	// AddConfiguration declares a REQUIRED configuration key of the given type. Like
	// AddStringConfiguration it sets required=true even though it accepts a defaultValue
	// (laatooserver core/configurableobject.go:178-181), so the default serves only as
	// documentation.
	AddConfiguration(ctx ServerContext, name string, desc string, dtype datatypes.DataType, defaultValue interface{})
	// AddOptionalConfiguration declares an optional configuration key. This is the only
	// declaration form whose default is actually reachable -- though it still comes back from the
	// Get* accessors with ok == false. Unlike AddStringConfiguration, desc is kept.
	AddOptionalConfiguration(ctx ServerContext, name string, desc string, dtype datatypes.DataType, defaultValue interface{})
	// GetConfiguration returns the configured value, falling back to the declared default. It
	// returns false BOTH for an undeclared key (value nil) and for a declared key that fell back
	// to its default (value non-nil); the two cases are distinguishable only by the value.
	GetConfiguration(ctx ServerContext, name string) (interface{}, bool)
	// GetStringConfiguration returns the value as a string. It PANICS when the value is present
	// but is not a string: the implementation asserts c.(string) with no comma-ok whenever the
	// value is non-nil (laatooserver core/configurableobject.go:207-213). Reading a key declared
	// as stringsmap, bool or config through this method crashes rather than returning false.
	GetStringConfiguration(ctx ServerContext, name string) (string, bool)
	// GetSecretConfiguration does NOT read this object's configurations at all. It ignores the
	// declaration map entirely and asks the SecretsManager for a secret literally named `name`
	// (laatooserver core/configurableobject.go:271-278) -- declaring `name` as a configuration is
	// neither required nor sufficient. It also type-asserts the SecretsManager server element
	// without a comma-ok, so it panics when no SecretsManager is registered, which makes the nil
	// check on the following line unreachable.
	GetSecretConfiguration(ctx ServerContext, name string) ([]byte, bool, error)
	// GetStringsMapConfiguration ALWAYS returns (nil, false) for a key declared
	// `type: stringsmap`. Configuration parsing stores that value as utils.StringsMap, a defined
	// type, while this method asserts the unnamed map[string]string
	// (laatooserver core/configurableobject.go:229-239); Go type assertions require identical
	// types, so the assertion can never succeed. Read such a key with GetConfiguration and assert
	// utils.StringsMap yourself.
	GetStringsMapConfiguration(ctx ServerContext, name string) (map[string]string, bool)
	// GetStringArrayConfiguration returns the value of a key declared `type: stringarr`. It
	// returns (nil, false) when the value is absent or is not a []string.
	GetStringArrayConfiguration(ctx ServerContext, name string) ([]string, bool)
	// GetBoolConfiguration returns the value as a bool. Like GetStringConfiguration it asserts
	// without a comma-ok and PANICS on a present, non-bool value
	// (laatooserver core/configurableobject.go:242-248) -- including a YAML default written as the
	// string "false".
	GetBoolConfiguration(ctx ServerContext, name string) (bool, bool)
	// GetMapConfiguration returns the value of a key declared `type: config` as a Config, and also
	// converts a raw map[string]interface{} into one. It does NOT work for a key declared
	// `type: stringmap`: parsing stores that as utils.StringMap, a defined type, and neither
	// assertion in the implementation matches it
	// (laatooserver core/configurableobject.go:251-268), so the result is (nil, false).
	GetMapConfiguration(ctx ServerContext, name string) (config.Config, bool)
}
