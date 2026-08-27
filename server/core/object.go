package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
)

// Creates object
type ObjectCreator func(ctx.Context) interface{}

// Creates collection
type ObjectCollectionCreator func(cx ctx.Context, length int) interface{}

// interface that needs to be implemented by any object provider in a system
type ObjectFactory interface {
	//Creates object
	CreateObject(ctx.Context) interface{}
	//Creates collection
	CreateObjectCollection(cx ctx.Context, length int) interface{}
	//Creates collection of pointers to object
	CreateObjectPointersCollection(cx ctx.Context, length int) interface{}
	//Get Metadata for the object
	Info() Info
}

// Initializable is the hook the server uses to configure a plain registered object -- one created
// through ctx.CreateObject rather than loaded as a service, factory or module. It is how a rule
// object receives its rule configuration and how the security handler seeds the Anonymous and
// System user objects.
//
// It is opt-in but NOT optional in practice: two of the server's three call sites type-assert to
// it WITHOUT a comma-ok (laatooserver core/rulesmanager.go:108 and
// core/securityhandler.go:213 and :228), so a rule object or a configured user object that does
// not implement Initialize panics at startup rather than being skipped or reported.
type Initializable interface {
	// Initialize receives the object's configuration. For a rule, conf is the rule YAML; for a
	// user object it is a synthesised config carrying Id and Roles. Returning an error aborts the
	// load at the rules-manager and security-handler call sites.
	Initialize(ctx ctx.Context, conf config.Config) error
}
